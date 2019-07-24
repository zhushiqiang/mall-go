package logic

import (
	"github.com/streadway/amqp"
	"github.com/yakaa/log4g"
	"log"
	"mall-go/integral/model"
	"mall-go/integral/protos"
	"context"
)

type (
	Integrallogic struct {
		dialHost string
		queueName string
		rabbitMqConn *amqp.Connection
		integralModel *model.IntegralModel
	}
)

func NewIntegralLogic(dataSource, QueueName string, integralModel *model.IntegralModel) (*Integrallogic,error)  {
	integrallogic := &Integrallogic{dialHost:dataSource,queueName:QueueName, integralModel:integralModel}
	if err :=integrallogic.CreateDial(); err != nil{
		return nil,err
	}
	return integrallogic,nil
}

func (l *Integrallogic) CreateDial() error  {
	conn, err := amqp.Dial(l.dialHost + l.queueName)

	if err != nil {
		return err
	}
	l.rabbitMqConn = conn
	return nil
}

func (l *Integrallogic) CloseRabbitMqConn()  {
	if err:= l.rabbitMqConn.Close();err != nil {
		log4g.ErrorFormat("CloseRabbitMqConn err %+v", err)
	}
}

func (l *Integrallogic) PushMessage(message string)  {
	ch, err := l.rabbitMqConn.Channel()
	defer func() {
		if err:= ch.Close();err != nil {
			log4g.ErrorFormat("CloseChConn err %+v", err)
		}
	}()
	q, err := l.QueueDeclare(ch)

	if err != nil {
		log4g.ErrorFormat("PushMessage err %+v", err)
		return
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing {
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (l *Integrallogic) ConsumeMessage(message string)  {
	ch, err := l.rabbitMqConn.Channel()
	defer func() {
		if err:= ch.Close();err != nil {
			log4g.ErrorFormat("CloseChConn err %+v", err)
		}
	}()
	if err !=nil {
		log4g.ErrorFormat("ConsumeMessage err %+v", err)
	}
	q, err := l.QueueDeclare(ch)

	if err != nil {
		log4g.ErrorFormat("ConsumeMessage err %+v", err)
		return
	}
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	if err !=nil {
		log4g.ErrorFormat("ch.Consume err %+v", err)
	}
	go func() {
		for d := range msgs  {
			msg := d.Body;
			//log.Printf("Received a message: %s", d.Body)
			if err := l.integralModel.ExecSql(string(msg)); err != nil {
				l.PushMessage(string(d.Body))
			}else {
				log4g.InfoFormat("Consume message %s [success]", msg)
			}
		}
	}()
}

func (l *Integrallogic) QueueDeclare(ch *amqp.Channel) (amqp.Queue,error)  {
	return ch.QueueDeclare(
		l.queueName, // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
}
func (l *Integrallogic) AddIntegral(_ context.Context, r *protos.AddIntegralRequest) (*protos.IntegralResponse, error) {

	l.PushMessage(l.integralModel.InsertIntegralSql(r.UserId, r.Integral))

	return &protos.IntegralResponse{
		UserId:r.UserId,Integral:r.Integral,
	},nil
}

func (l *Integrallogic) ConsumerIntegral(_ context.Context, r *protos.ConsumerIntegralRequest) (*protos.IntegralResponse, error) {
	l.PushMessage(l.integralModel.UpdateIntegralByUserIdSql(r.UserId, r.ConsumerIntegral))

	return new(protos.IntegralResponse),nil
}

func (l *Integrallogic) FindOneByUserId(_ context.Context, r *protos.FindOneByUserIdRequest) (*protos.IntegralResponse, error) {
	integral, err := l.integralModel.FindByUserId(r.UserId)

	if err != nil {
		return nil,err
	}

	return &protos.IntegralResponse{
		UserId:integral.UserId,Integral:integral.Integral,
	},nil
}
