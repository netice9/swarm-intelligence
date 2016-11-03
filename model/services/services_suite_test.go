package services_test

import (
	"github.com/docker/docker/api/types/swarm"
	"github.com/netice9/swarm-intelligence/model"
	"github.com/netice9/swarm-intelligence/model/services"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestServices(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}

type RecordedEvent struct {
	Event interface{}
	Args  []interface{}
}

type FakeEmitter struct {
	Events []RecordedEvent
}

func (f *FakeEmitter) Emit(event interface{}, arguments ...interface{}) model.EventEmitter {
	f.Events = append(f.Events, RecordedEvent{event, arguments})
	return f
}
func (f *FakeEmitter) RemoveListener(event, listener interface{}) model.EventEmitter { return f }
func (f *FakeEmitter) AddListener(event, listener interface{}) model.EventEmitter    { return f }

var _ = Describe("ServicesAggregator", func() {

	var aggregator *services.ServicesAggregator
	var emitter *FakeEmitter

	BeforeEach(func() {
		emitter = &FakeEmitter{}
		aggregator = services.NewServicesAggregator(emitter)
	})

	Describe("OnServices()", func() {
		var update []swarm.Service

		Context("When aggregator haven't received any updates", func() {
			Context("When event contains one new service", func() {

				BeforeEach(func() {
					update = []swarm.Service{
						swarm.Service{
							ID: "id1",
							Spec: swarm.ServiceSpec{
								Annotations: swarm.Annotations{
									Name: "service1",
								},
							},
						},
					}

				})

				BeforeEach(func() {
					aggregator.OnServices(update)
				})

				It("Should add service name and ID to the list of services", func() {
					Expect(aggregator.ServiceList()).To(Equal([]services.ServiceStatus{{Name: "service1", ID: "id1"}}))
				})

				It("Should fire new 'list' event with the list of services and new services update event", func() {
					Expect(emitter.Events).To(Equal(
						[]RecordedEvent{
							{"update/id1", []interface{}{update[0]}},
							{"list", []interface{}{services.ServiceList{{Name: "service1", ID: "id1"}}}},
						}))
				})
				Context("When same event is received again", func() {
					BeforeEach(func() {
						aggregator.OnServices(update)
					})

					It("Should fire new 'list' event with the list of services and new services update event only once", func() {
						Expect(emitter.Events).To(Equal(
							[]RecordedEvent{
								{"update/id1", []interface{}{update[0]}},
								{"list", []interface{}{services.ServiceList{{Name: "service1", ID: "id1"}}}},
							}))
					})

				})

				Context("When aggregator received empty list of events", func() {
					BeforeEach(func() {
						update = []swarm.Service{}
						aggregator.OnServices(update)
					})

					It("Should remove service name from the list of services", func() {
						Expect(aggregator.ServiceList()).To(Equal([]services.ServiceStatus{}))
					})

					It("Should fire new 'delete' event for the deleted service", func() {
						Expect(emitter.Events[2]).To(Equal(RecordedEvent{Event: "delete/id1"}))
					})

					It("Should fire new 'list' event with the empty list of services", func() {
						Expect(emitter.Events[3]).To(Equal(RecordedEvent{Event: "list", Args: []interface{}{services.ServiceList{}}}))
					})

				})

			})
		})
	})
})
