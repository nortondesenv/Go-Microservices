package kafka

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/nortondesenv/Go-Microservice/internal/models"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"github.com/segmentio/kafka-go"
)

func (pcg *ProductsConsumerGroup) createProductWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductsConsumerGroup.createProductWorker")
	defer span.Finish()

	span.LogFields(log.String("ConsumerGroup", r.Config().GroupID))

	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			pcg.log.Errorf("FetchMessage", err)
			return
		}

		pcg.log.Infof(
			"WORKER: %v, message at topic/partition/offset %v/%v/%v: %s = %s\n",
			workerID,
			m.Topic,
			m.Partition,
			m.Offset,
			string(m.Key),
			string(m.Value),
		)

		var prod models.Product
		if err := json.Unmarshal(m.Value, &prod); err != nil {
			pcg.log.Errorf("json.Unmarshal", err)
			continue
		}

		created, err := pcg.productsUC.Create(ctx, &prod)
		if err != nil {
			pcg.log.Errorf("productsUC.Create", err)
			continue
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			pcg.log.Errorf("FetchMessage", err)
			continue
		}

		pcg.log.Infof("Created product: %v", created)
	}
}

func (pcg *ProductsConsumerGroup) updateProductWorker(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ProductsConsumerGroup.updateProductWorker")
	defer span.Finish()

	span.LogFields(log.String("ConsumerGroup", r.Config().GroupID))

	defer wg.Done()

	for {
		m, err := r.FetchMessage(ctx)
		if err != nil {
			pcg.log.Errorf("FetchMessage", err)
			return
		}

		pcg.log.Infof(
			"WORKER: %v, message at topic/partition/offset %v/%v/%v: %s = %s\n",
			workerID,
			m.Topic,
			m.Partition,
			m.Offset,
			string(m.Key),
			string(m.Value),
		)

		var prod models.Product
		if err := json.Unmarshal(m.Value, &prod); err != nil {
			pcg.log.Errorf("json.Unmarshal", err)
			continue
		}

		updated, err := pcg.productsUC.Update(ctx, &prod)
		if err != nil {
			pcg.log.Errorf("productsUC.Update", err)
			continue
		}

		if err := r.CommitMessages(ctx, m); err != nil {
			pcg.log.Errorf("FetchMessage", err)
			continue
		}

		pcg.log.Infof("Created product: %v", updated)
	}
}
