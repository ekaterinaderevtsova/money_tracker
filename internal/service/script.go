package service

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type ScriptService struct {
	spendingRepository ISpendingRepository
}

func NewScriptService(spendingRepository ISpendingRepository) *ScriptService {
	return &ScriptService{spendingRepository: spendingRepository}
}

func (ss *ScriptService) TransferToRedis(ctx context.Context) error {
	date := time.Now()

	spendings, err := ss.spendingRepository.GetWeekSpendings(ctx, date)
	if err != nil {
		return err
	}

	for _, daySpending := range spendings.DaySpendings {
		if daySpending.Sum != 0 {
			parts := strings.Split(daySpending.Day, "-")
			daySpending.Day = fmt.Sprintf("2025-%s-%s", parts[1], parts[0])

			err = ss.spendingRepository.AddCurrentWeekSpending(ctx, &daySpending)
			if err != nil {
				return err
			}
		}

	}

	return nil
}
