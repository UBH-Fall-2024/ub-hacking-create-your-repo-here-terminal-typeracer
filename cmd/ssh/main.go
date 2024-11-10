package main

// An example Bubble Tea server. This will put an ssh session into alt screen
// and continually print up to date terminal information.

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fejiberglibstein/terminal-typeracer/pkg/network"
	"github.com/Fejiberglibstein/terminal-typeracer/pkg/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
)

const (
	host = "0.0.0.0"
	port = "23234"
)

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbleteaMiddleware(),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}

func bubbleteaMiddleware() wish.Middleware {
	newProg := func(m *tui.Model, opts ...tea.ProgramOption) *tea.Program {
		p := tea.NewProgram(m, opts...)

		m.SendMessage(&network.Message{
			Header: uint8(network.Connect),
			Data:   m.Sess.User(),
		})

		// Handle all events sent by the server
		go func() {
			for {
				message, err := m.ReadMessage()
				if err != nil {
					log.Print("Could not read from server: ", err)
					continue
				}
				p.Send(message)
			}
		}()

		return p
	}
	teaHandler := func(s ssh.Session) *tea.Program {
		pty, _, active := s.Pty()
		renderer := bubbletea.MakeRenderer(s)
		if !active {
			wish.Fatalln(s, "no active terminal, skipping")
			return nil
		}

		conn, err := net.Dial("tcp", net.JoinHostPort(host, port))
		if err != nil {
			log.Print("Could not connect the poor soul :(")
			return nil
		}

		m := tui.NewModel(renderer, &pty, &conn, s)
		return newProg(m, append(bubbletea.MakeOptions(s), tea.WithAltScreen())...)
	}
	return bubbletea.MiddlewareWithProgramHandler(teaHandler, termenv.ANSI256)
}
