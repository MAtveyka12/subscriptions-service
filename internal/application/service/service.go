package service

type Service interface {
	SubscriptionService
}

type service struct {
	SubscriptionService
}

func NewService(subscriptionService SubscriptionService) Service {
	return &service{
		SubscriptionService: subscriptionService,
	}
}
