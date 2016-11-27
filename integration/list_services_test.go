package integration_test

import (
	"log"
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
		BeforeEach(func() {
			time.Sleep(500 * time.Millisecond)
			Expect(page.Navigate("http://localhost:5080")).To(Succeed())
		})
		It("Should list the x service", func() {
			Eventually(page.First("div.panel-heading")).Should(BeFound())
			html, err := page.HTML()
			Expect(err).ToNot(HaveOccurred())
			log.Println(html)
		})
	})
})
