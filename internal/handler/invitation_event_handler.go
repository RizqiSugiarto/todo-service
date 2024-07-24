package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/digisata/invitation-service/internal/constant"
	"github.com/digisata/invitation-service/internal/entity"
	"github.com/digisata/invitation-service/pkg/rabbitclient"
	amqp "github.com/rabbitmq/amqp091-go"
)

type InvitationEventHandler struct {
	invitationUseCase InvitationUseCase
	rabbitMQClient    *rabbitclient.RabbitMQ
	consumerCount     int
}

func NewInvitationEvent(invitationUseCase InvitationUseCase, rabbitMQclient *rabbitclient.RabbitMQ, consumerCount int) *InvitationEventHandler {
	return &InvitationEventHandler{
		invitationUseCase: invitationUseCase,
		rabbitMQClient:    rabbitMQclient,
		consumerCount:     consumerCount,
	}
}

func (h *InvitationEventHandler) StartConsumeInvitationOpenEvent() error {
	err := h.rabbitMQClient.StartConsuming(constant.INVITATION_OPEN_QUEUE, h.consumerCount, func(d amqp.Delivery) {
		var event entity.InvitationEventRequest
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			log.Println("Failed to unmarshal message:", err.Error())
			return
		}

		isOpen := true
		req := entity.UpdateInvitationRequest{
			ID:     event.ID,
			IsOpen: &isOpen,
			OpenAt: event.OpenAt,
		}

		err = h.invitationUseCase.UpdateInvitation(context.Background(), req)
		if err != nil {
			log.Println("Failed to update invitation:", err.Error())
			return
		}
	})
	if err != nil {
		return fmt.Errorf("failed to start RabbitMQ consumers for %s queue: %s", constant.INVITATION_OPEN_QUEUE, err.Error())
	}

	return nil
}

func (h *InvitationEventHandler) StartConsumeInvitationComingEvent() error {
	err := h.rabbitMQClient.StartConsuming(constant.INVITATION_COMING_QUEUE, h.consumerCount, func(d amqp.Delivery) {
		var event entity.InvitationEventRequest
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			log.Println("Failed to unmarshal message:", err.Error())
			return
		}

		req := entity.UpdateInvitationRequest{
			ID:       event.ID,
			IsComing: event.IsComing,
		}

		err = h.invitationUseCase.UpdateInvitation(context.Background(), req)
		if err != nil {
			log.Println("Failed to update invitation:", err.Error())
			return
		}
	})
	if err != nil {
		return fmt.Errorf("failed to start RabbitMQ consumers for %s queue: %s", constant.INVITATION_COMING_QUEUE, err.Error())
	}

	return nil
}

func (h *InvitationEventHandler) StartConsumeInvitationSendMoneyEvent() error {
	err := h.rabbitMQClient.StartConsuming(constant.INVITATION_SEND_MONEY_QUEUE, h.consumerCount, func(d amqp.Delivery) {
		var event entity.InvitationEventRequest
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			log.Println("Failed to unmarshal message:", err.Error())
			return
		}

		isSendMoney := true
		req := entity.UpdateInvitationRequest{
			ID:          event.ID,
			IsSendMoney: &isSendMoney,
		}

		err = h.invitationUseCase.UpdateInvitation(context.Background(), req)
		if err != nil {
			log.Println("Failed to update invitation:", err.Error())
			return
		}
	})
	if err != nil {
		return fmt.Errorf("failed to start RabbitMQ consumers for %s queue: %s", constant.INVITATION_SEND_MONEY_QUEUE, err.Error())
	}

	return nil
}

func (h *InvitationEventHandler) StartConsumeInvitationSendGiftEvent() error {
	err := h.rabbitMQClient.StartConsuming(constant.INVITATION_SEND_GIFT_QUEUE, h.consumerCount, func(d amqp.Delivery) {
		var event entity.InvitationEventRequest
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			log.Println("Failed to unmarshal message:", err.Error())
			return
		}

		isSendGift := true
		req := entity.UpdateInvitationRequest{
			ID:         event.ID,
			IsSendGift: &isSendGift,
		}

		err = h.invitationUseCase.UpdateInvitation(context.Background(), req)
		if err != nil {
			log.Println("Failed to update invitation:", err.Error())
			return
		}
	})
	if err != nil {
		return fmt.Errorf("failed to start RabbitMQ consumers for %s queue: %s", constant.INVITATION_SEND_GIFT_QUEUE, err.Error())
	}

	return nil
}

func (h *InvitationEventHandler) StartConsumeInvitationCheckInEvent() error {
	err := h.rabbitMQClient.StartConsuming(constant.INVITATION_CHECK_IN_QUEUE, h.consumerCount, func(d amqp.Delivery) {
		var event entity.InvitationEventRequest
		err := json.Unmarshal(d.Body, &event)
		if err != nil {
			log.Println("Failed to unmarshal message:", err.Error())
			return
		}

		isCheckIn := true
		req := entity.UpdateInvitationRequest{
			ID:        event.ID,
			IsCheckIn: &isCheckIn,
			CheckInAt: event.CheckInAt,
		}

		err = h.invitationUseCase.UpdateInvitation(context.Background(), req)
		if err != nil {
			log.Println("Failed to update invitation:", err.Error())
			return
		}
	})
	if err != nil {
		return fmt.Errorf("failed to start RabbitMQ consumers for %s queue: %s", constant.INVITATION_CHECK_IN_QUEUE, err.Error())
	}

	return nil
}
