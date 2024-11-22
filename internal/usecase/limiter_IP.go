package usecase

import "github.com/douglasschantz/rate-limiter/internal/entity"

type RequestIPInputDTO struct {
	IP string `json:"ip"`
}

type CreateLimiterIPUseCase struct {
	LimiterIPRepository entity.RequestRepositoryInterface
}

func NewCreateLimiterIPUseCase(LimiterIPRepository entity.RequestRepositoryInterface) *CreateLimiterIPUseCase {
	return &CreateLimiterIPUseCase{LimiterIPRepository: LimiterIPRepository}
}

func (c *CreateLimiterIPUseCase) Execute(input RequestIPInputDTO) error {
	respIp := entity.Ip{
		IP: input.IP,
	}

	if err := c.LimiterIPRepository.SaveRequestIP(&respIp); err != nil {
		return err
	}
	resultCount, err := c.LimiterIPRepository.GetCountLimiter(&respIp)
	if err != nil {
		return nil
	}
	count := entity.Count{
		Count: resultCount,
	}
	count.CheckLimiterIp()
	return nil
}
