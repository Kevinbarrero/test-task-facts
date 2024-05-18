package launcher

import (
	"log"
	"sync"
	"test-task-facts/external"
	"test-task-facts/utils"
)

func RunOps(nFacts int) {
	// Это выполняемая операция, она используется только для того, чтобы узнать, в каком состоянии находится лог.
	op := "RunOps"
	// канал для фактов сохранен правильно
	saveFactChan := make(chan external.Fact)
	// канал для ошибок на случай сохранения факта
	errChan := make(chan error)
	var wg sync.WaitGroup

	// Запускаем горутины nFacts
	for i := 0; i < nFacts; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// utils.generateFact создает случайный факт
			fact := utils.GenerateFact()
			// external.SaveFact сохраняет факт в базу данных
			factResp, err := external.SaveFact(fact)
			if err != nil {
				// В случае ошибки отправляем сигнал в errChan и завершаем горутину
				errChan <- err
				return
			}
			// В случае успеха отправляем сигнал в saveFactChan
			saveFactChan <- factResp
		}()
	}

	// закрыть каналы и дождаться wg в конце выполнения
	go func() {
		wg.Wait()
		close(saveFactChan)
		close(errChan)
		log.Println(op, nFacts, "Facts saved correctly")
	}()

	// в момент поступления сигнала на factResp отправляем запрос на external.IsFactsaved,
	// чтобы проверить, действительно ли он был сохранен.
	go func() {
		for factResp := range saveFactChan {
			wg.Add(1)
			factResp := factResp
			go func() {
				defer wg.Done()
				isSaved, err := external.IsFactSaved(&factResp)
				if err != nil {
					log.Println(op, err)
					return
				}
				log.Println("double check with indicators/get_facts, Fact ID:", factResp.IndicatorToMoFactID, "is saved:", isSaved)
			}()
		}
	}()

	// в момент поступления сигнала на errChan пытаемся сохранить факт еще раз.
	// Если не удалось, снова отправляем сигнал в errChan.
	go func() {
		for err := range errChan {
			wg.Add(1)
			err := err
			go func() {
				log.Println(op, err, "trying to save again...")
				defer wg.Done()
				fact := utils.GenerateFact()
				factResp, ferr := external.SaveFact(fact)
				if ferr != nil {
					errChan <- ferr
					return
				}
				saveFactChan <- factResp
			}()
		}
	}()
	wg.Wait()
}
