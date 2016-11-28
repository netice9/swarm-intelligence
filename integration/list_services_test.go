package integration_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sclevine/agouti"
	. "github.com/sclevine/agouti/matchers"
)

var _ = Describe("List services", func() {
	var page *agouti.Page

	BeforeEach(func() {
		var err error
		page, err = agoutiDriver.NewPage()
		Expect(err).NotTo(HaveOccurred())
	})

	Context("When I visit the Index page", func() {
		BeforeEach(func(done Done) {
			defer close(done)
			for {
				Expect(page.Navigate("http://localhost:5080")).To(Succeed())
				selection := page.FindByID("react-application")
				count, _ := selection.Count()
				if count == 1 {
					return
				}
				time.Sleep(100 * time.Millisecond)
			}
		}, 3.0)
		It("Should list the 'swarm-intelligence-head' service", func() {
			Eventually(page.First("h4.list-group-item-heading")).Should(HaveText("swarm-intelligence-head"))
		})
	})
})
