package hcsprocess_test

import (
	"github.com/genevieve/hcsprocess"
	"github.com/genevieve/hcsprocess/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Process", func() {
	var wrappedProcess *hcsprocess.Process
	var fakeProcess *fakes.Process

	BeforeEach(func() {
		fakeProcess = &fakes.Process{}
		wrappedProcess = hcsprocess.New(fakeProcess)
	})

	Describe("AttachIO", func() {
		var (
			attachedStdin,
			attachedStdout,
			attachedStderr,
			processStdin,
			processStdout,
			processStderr *gbytes.Buffer
		)

		BeforeEach(func() {
			attachedStdin = gbytes.BufferWithBytes([]byte("something-on-stdin"))
			attachedStdout = gbytes.NewBuffer()
			attachedStderr = gbytes.NewBuffer()

			processStdin = gbytes.NewBuffer()
			processStdout = gbytes.BufferWithBytes([]byte("something-on-stdout"))
			processStderr = gbytes.BufferWithBytes([]byte("something-on-stderr"))

			fakeProcess.StdioCall.Returns.Stdin = processStdin
			fakeProcess.StdioCall.Returns.Stdout = processStdout
			fakeProcess.StdioCall.Returns.Stderr = processStderr
		})

		It("attaches process IO to stdin, stdout, and stderr", func() {
			exitCode, err := wrappedProcess.AttachIO(attachedStdin, attachedStdout, attachedStderr)
			Expect(err).NotTo(HaveOccurred())

			Expect(exitCode).To(Equal(0))

			Eventually(processStdin).Should(gbytes.Say("something-on-stdin"))
			Eventually(attachedStdout).Should(gbytes.Say("something-on-stdout"))
			Eventually(attachedStderr).Should(gbytes.Say("something-on-stderr"))
		})

		It("closes the process' stdin pipe after copying", func() {
			exitCode, err := wrappedProcess.AttachIO(attachedStdin, nil, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(exitCode).To(Equal(0))

			Eventually(processStdin).Should(gbytes.Say("something-on-stdin"))
			Expect(fakeProcess.CloseStdinCall.CallCount).To(Equal(1))
		})

		Context("when stdout and stderr are done, but attached stdin is blocking", func() {
			var neverendingAttachedStdin *fakes.Reader

			BeforeEach(func() {
				neverendingAttachedStdin = &fakes.Reader{}
			})

			It("should exit before the 5 second timeout", func() {
				code := make(chan int)
				go func() {
					exitCode, err := wrappedProcess.AttachIO(neverendingAttachedStdin, attachedStdout, attachedStderr)
					Expect(err).NotTo(HaveOccurred())
					code <- exitCode
				}()

				// AttachIO should exit within 1 second because stdout & sdtderr are done.
				Eventually(code).Should(Receive(Equal(0), "AttachIO didn't exit."))

				Expect(attachedStdout).To(gbytes.Say("something-on-stdout"))
				Expect(attachedStderr).To(gbytes.Say("something-on-stderr"))
			})
		})
	})
})
