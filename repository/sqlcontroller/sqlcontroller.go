package sqlcontroller

import (
	"core-bot/config"
	"core-bot/controllers/botcontroller"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func SqlCommands(bot *tgbotapi.BotAPI, updates tgbotapi.UpdatesChannel, userState map[int64]string) {
	// Variabel global untuk menyimpan state query saat ini (DDMAST, CDMAST, LNMAST)
	currentQueryState := make(map[int64]string) // Map to store each chat's current query state

	for update := range updates {
		if update.Message != nil { // If we receive a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			chatID := update.Message.Chat.ID
			state, exists := userState[chatID]

			if !exists {
				// Jika user tidak memiliki state, kirim pesan salam
				botcontroller.SendGreeting(bot, update, userState)
				userState[chatID] = "awaiting_command"
			} else {
				switch state {
				case "awaiting_command":
					switch update.Message.Text {
					case "halo":
						botcontroller.SendGreeting(bot, update, userState)
					case "1":
						// Ubah state menjadi 'awaiting_data_1'
						userState[chatID] = "awaiting_data"
						currentQueryState[chatID] = "DDMAST" // Set current query
						msg := tgbotapi.NewMessage(chatID, "Silakan masukkan data untuk query DDMAST seperti ini:\n\nStatus: 1\nTipe Rek: S\nMata uang: IDR\n\nSetiap field dapat kamu ganti sesuai dengan kebutuhannya. Jika ingin kembali ke menu utama, ketik 'kembali'.")
						bot.Send(msg)
					case "2":
						// Ubah state menjadi 'awaiting_data_2'
						userState[chatID] = "awaiting_data"
						currentQueryState[chatID] = "CDMAST" // Set current query
						msg := tgbotapi.NewMessage(chatID, "Silakan masukkan data untuk query CDMAST seperti ini:\n\nStatus:1\nMata uang: IDR\n\nSetiap field dapat kamu ganti sesuai dengan kebutuhannya. Jika ingin kembali ke menu utama, ketik 'kembali'.")
						bot.Send(msg)
					case "3":
						// Ubah state menjadi 'awaiting_data_3'
						userState[chatID] = "awaiting_data"
						currentQueryState[chatID] = "LNMAST" // Set current query
						msg := tgbotapi.NewMessage(chatID, "Silakan masukkan data untuk query LNMAST seperti ini:\n\nStatus:1\nMata uang: IDR\n\nSetiap field dapat kamu ganti sesuai dengan kebutuhannya. Jika ingin kembali ke menu utama, ketik 'kembali'.")
						bot.Send(msg)
					default:
						msg := tgbotapi.NewMessage(chatID, "Perintah tidak dikenali. Silakan pilih opsi:\n 1. Query DDMAST \n 2. Query CDMAST \n 3. Query LNMAST")
						bot.Send(msg)
					}

				case "awaiting_data":
					// Check if the user wants to go back
					if strings.Contains(strings.ToLower(update.Message.Text), "kembali") {
						userState[chatID] = "awaiting_command"
						msg := tgbotapi.NewMessage(chatID, "Anda telah kembali ke menu utama. Silakan pilih opsi:\n 1. Query DDMAST \n 2. Query CDMAST \n 3. Query LNMAST")
						bot.Send(msg)
						continue // Keep listening for further commands
					}

					query := currentQueryState[chatID]
					handleQuery(bot, update, query, userState)
					// Setelah query selesai, ubah state ke 'awaiting_continue_query'
					userState[chatID] = "awaiting_continue_query"
					msg := tgbotapi.NewMessage(chatID, "Apakah Anda ingin melanjutkan query? (ya/tidak)")
					bot.Send(msg)

				case "awaiting_continue_query":
					userMessage := strings.ToLower(update.Message.Text)

					if userMessage == "ya" {
						// Kembali ke state 'awaiting_data' sesuai query yang sebelumnya
						userState[chatID] = "awaiting_data"
						msg := tgbotapi.NewMessage(chatID, "Silakan masukkan data untuk query berikutnya.")
						bot.Send(msg)
					} else {
						// Kembali ke state 'awaiting_command'
						userState[chatID] = "awaiting_command"
						msg := tgbotapi.NewMessage(chatID, "Anda telah kembali ke menu utama. Silakan pilih opsi:\n 1. Query DDMAST \n 2. Query CDMAST \n 3. Query LNMAST")
						bot.Send(msg)
					}

				default:
					msg := tgbotapi.NewMessage(chatID, "Perintah tidak dikenali.")
					bot.Send(msg)
				}
			}
		}
	}
}

func handleQuery(bot *tgbotapi.BotAPI, update tgbotapi.Update, tableName string, userState map[int64]string) {
	chatID := update.Message.Chat.ID
	userMessage := update.Message.Text

	// Variabel untuk menyimpan hasil parsing
	var status, actype, currency, sccode, cbal, npdt6, limitdt string

	// Bagi input menjadi baris-baris
	inputLines := strings.Split(userMessage, "\n")

	var validateStatus, validateActype, validateCurrency bool

	// Process each line to extract values
	for _, line := range inputLines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Status") {
			status = strings.TrimSpace(strings.TrimPrefix(line, "Status"))
			validateStatus = true
		} else if strings.HasPrefix(line, "Tipe") {
			actype = strings.TrimSpace(strings.TrimPrefix(line, "Tipe"))
			validateActype = true
		} else if strings.HasPrefix(line, "Mata uang") {
			currency = strings.TrimSpace(strings.TrimPrefix(line, "Mata uang"))
			validateCurrency = true
		} else if strings.HasPrefix(line, "SCCODE") {
			sccode = strings.TrimSpace(strings.TrimPrefix(line, "SCCODE"))
		} else if strings.HasPrefix(line, "CBAL") {
			cbal = strings.TrimSpace(strings.TrimPrefix(line, "CBAL"))
		} else if strings.HasPrefix(line, "NPDT6") {
			npdt6 = strings.TrimSpace(strings.TrimPrefix(line, "NPDT6"))
		} else if strings.HasPrefix(line, "Limit") {
			limitdt = strings.TrimSpace(strings.TrimPrefix(line, "Limit"))
		}
	}

	// Validate required fields
	if !validateStatus || !validateActype || !validateCurrency {
		// Change state to 'awaiting_command'
		userState[chatID] = "awaiting_command"

		// Send message indicating invalid input and return to main menu
		msg := tgbotapi.NewMessage(chatID, "Input tidak sesuai. Anda telah kembali ke menu utama. Silakan pilih opsi berikut:\n 1. Query DDMAST \n 2. Query CDMAST \n 3. Query LNMAST")
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Failed to send message:", err)
		}
		return
	}

	// Log extracted values
	log.Printf("Status: %s, Actype: %s, Mata Uang: %s, SCCODE: %s", status, actype, currency, sccode)

	// Generate SQL Query
	sqlQuery := fmt.Sprintf("SELECT ACCTNO FROM %s WHERE STATUS = %s AND ACTYPE = '%s' AND DDCTYP = '%s'", tableName, status, actype, currency)

	if sccode != "" {
		// If SCCODE is provided, add to query
		sqlQuery += fmt.Sprintf(" AND SCCODE = '%s'", sccode)
	}
	if cbal != "" {
		sqlQuery += fmt.Sprintf(" AND CBAL = '%s'", cbal)
	}
	if npdt6 != "" {
		sqlQuery += fmt.Sprintf(" AND NPDT6 = '%s'", npdt6)
	}
	if limitdt != "" {
		// If LIMIT is provided, add to query
		sqlQuery += fmt.Sprintf(" LIMIT %s", limitdt)

	} else {
		// If LIMIT is not provided, use default value
		sqlQuery += " LIMIT 10"
	}

	log.Printf("Generated SQL Query: %s", sqlQuery)

	// Ambil hasil query
	db := config.GetDB()
	rows, err := db.Query(sqlQuery)
	if err != nil {
		// Handle SQL error and inform the user
		msg := tgbotapi.NewMessage(chatID, "Perintah tidak dikenali. Silakan cek kembali inputan Anda dan coba lagi.")
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Failed to send message:", err)
		}
		// Reset state to continue
		userState[chatID] = "awaiting_command"
		return
	}
	defer rows.Close()

	var results []string
	for rows.Next() {
		var acctNo string
		if err := rows.Scan(&acctNo); err != nil {
			msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Error scanning results: %s", err))
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("Failed to send message:", err)
			}
			continue
		}
		results = append(results, acctNo)
	}

	if len(results) == 0 {
		results = append(results, "No results found.")
	}
	resultMessage := "Query Results:\n" + strings.Join(results, "\n")
	botcontroller.SendQueryResult(bot, update, resultMessage)
}
