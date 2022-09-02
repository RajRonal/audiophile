package Ginkgo

import (
	"audioPhile/handlers"
	"github.com/onsi/gomega/ghttp"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGinkgoCLI(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ginkgo CLI Suite")
}

var _ = Describe("Client", func() {
	var (
		server     *ghttp.Server
		statusCode int
		body       []byte
		path       string
		addr       string
	)
	BeforeEach(func() {
		// start a test http server
		server = ghttp.NewServer()
	})
	AfterEach(func() {
		server.Close()
	})

	Context("When Post request is sent to read path but file exists", func() {
		BeforeEach(func() {
			_, err := os.Create("data.txt")
			Expect(err).NotTo(HaveOccurred())
			//body := []byte(`{"first_name":"hello","last_name":"world","email":"hello@pqr.com","contact_number":"98989889","user_name":"helloworld","password":"helloworld"}`)
			//file.Write(body)
			body = []byte(`{"FirstName":"hello","LastName":"1234","Email":"hello@pqr.com","ContactNumber":"98989889","UserName":"helloworld","Password":"helloworld"}`)
			request := CreateRequest(body)

			//statusCode = 200
			path = "/api/signup"
			addr = "http://" + server.Addr() + path
			//url :=
			server.AppendHandlers(
				handlers.SignUp,
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", request.URL),
					ghttp.RespondWithPtr(&statusCode, &body),
				))

		})
		//AfterEach(func() {
		//	err := os.Remove("data.txt")
		//	Expect(err).NotTo(HaveOccurred())
		//})
		It("Reads data from file successfully", func() {
			bdy, err := getResponse(addr)
			Expect(err).ShouldNot(HaveOccurred())
			//Expect().Should(Equal(statusCode))
			Expect(bdy).To(Equal(body))
		})
	})
})
