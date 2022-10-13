package gatherer

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/miekg/dns"
)

type empty struct{}
type result struct {
	Ipaddress string
	Hostname  string
}

func lookupA(domain1, serveraddr string) ([]string, error) {
	var msg dns.Msg
	var ips []string
	fqdn := dns.Fqdn(domain1)
	msg.SetQuestion(fqdn, dns.TypeA)
	msg1, err := dns.Exchange(&msg, serveraddr)
	if err != nil {
		return ips, err
	}

	if len(msg1.Answer) < 1 {
		fmt.Println("[!] No records!")
		return ips, errors.New("no records")
	}

	for _, answer := range msg1.Answer {
		if msg_a, ok := answer.(*dns.A); ok {
			ips = append(ips, msg_a.A.String())
		}
	}

	return ips, nil
}

func lookupCNAME(domain1, serveraddr string) ([]string, error) {
	var msg dns.Msg
	var domainhosts []string

	fqdn := dns.Fqdn(domain1)
	msg.SetQuestion(fqdn, dns.TypeCNAME)
	msg_1, err := dns.Exchange(&msg, serveraddr)
	if err != nil {
		return domainhosts, err
	}

	if len(msg_1.Answer) < 1 {
		return domainhosts, errors.New("no answer")
	}

	for _, answer := range msg.Answer {
		if msg_a, ok := answer.(*dns.CNAME); ok {
			domainhosts = append(domainhosts, msg_a.Target)
		}
	}

	return domainhosts, nil
}

func lookup(domain, serveraddr string) []result {
	var results []result
	var domain_1 = domain

	for {
		cnames, err := lookupCNAME(domain_1, serveraddr)
		if err == nil && len(cnames) > 0 {
			domain_1 = cnames[0]
			continue
		}

		ips, err := lookupA(domain_1, serveraddr)
		if err != nil {
			break
		}
		for _, ip := range ips {
			results = append(results, result{Hostname: domain, Ipaddress: ip})
		}
		break
	}
	return results
}

func worker1(traker chan empty, domains chan string, result1 chan []result, serveraddr string) {
	for domain := range domains {
		results := lookup(domain, serveraddr)
		if len(results) > 0 {
			result1 <- results
		}
	}
	var ep1 empty
	traker <- ep1
}

func CollectingSubdomainNames(workers int, domain_base, wordlist, serveraddr string) {

	fmt.Println("[*] Collecting Subdomain Names .......")
	if domain_base == "" {
		fmt.Println("[!] Wordlist or Domain is required")
		os.Exit(1)
	}

	if wordlist == "" {
		fmt.Println("[*] 使用默认字典")
		wordlist = "txt\\SuperWordlist\\SubDomain.txt"
	}

	var results []result
	domains := make(chan string, workers)
	result1 := make(chan []result)
	traker := make(chan empty)

	file, err := os.Open(wordlist)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		domains <- fmt.Sprintf("%s.%s", scanner.Text(), domain_base)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner has a error", os.Stderr)
	}

	for i := 0; i < workers; i++ {
		go worker1(traker, domains, result1, serveraddr)
	}

	go func() {
		for r := range result1 {
			results = append(results, r...)
		}
		var e1 empty
		traker <- e1
	}()

	close(domains)

	for i := 0; i < workers; i++ {
		<-traker
	}
	close(result1)
	<-traker

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 4, ' ', 0)
	for _, rst := range results {
		fmt.Fprintf(w, "%s\t%s\n", rst.Hostname, rst.Ipaddress)
	}
	w.Flush()
}
