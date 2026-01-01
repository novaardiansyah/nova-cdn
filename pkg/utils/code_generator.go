package utils

import (
	"fmt"
	"math/rand"
	"time"

	"nova-cdn/internal/repositories"
)

func GetCode(repo *repositories.GenerateRepository, alias string, isNotPreview bool) string {
	gen, err := repo.FindByAlias(alias)
	if err != nil || gen == nil {
		return fmt.Sprintf("ER-%05d", rand.Intn(90000)+10000)
	}

	now := time.Now()
	date := now.Format("060102")
	separator := gen.Separator

	separatorTime, err := time.Parse("060102", separator)
	if err != nil {
		separatorTime = now
	}
	separatorStr := separatorTime.Format("060102")

	if gen.Queue == 9999 || date[:4] != separatorStr[:4] {
		gen.Queue = 1
		gen.Separator = date
	}

	queue := fmt.Sprintf("%s%04d%s", date[:4], gen.Queue, date[4:6])

	if gen.Prefix != nil && *gen.Prefix != "" {
		queue = *gen.Prefix + queue
	}
	if gen.Suffix != nil && *gen.Suffix != "" {
		queue = queue + *gen.Suffix
	}

	if isNotPreview {
		gen.Queue += 1
		repo.Update(gen)
	}

	return queue
}
