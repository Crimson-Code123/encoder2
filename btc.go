package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
)

/* Mining functionality */

const (
	pool    = "solo.ckpool.org:3333"
	address = ""
)

// // goroutine
// func init() {
// 	for {

// 	}
// }

func getHeight() int {
	resp, err := http.Get("https://blockchain.info/latestblock")
	if err != nil {
		fmt.Println("http:", err)
	}
	_ = resp
	// fmt.Println(resp.Body.Read())
	return 0
}

func miner() {

}

func listener() {
	conn, err := net.Dial("tcp", pool)
	if err != nil {
		fmt.Println("Listener:", err)
	}
	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)
	_ = r
	// wr := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	//https://go.dev/blog/strings
	init1 := `{"id": 1, "method": "mining.subscribe", "params": []}`
	init2 := `{"params": ["' + address.encode() + b'", "password"], "id": 2, "method": "mining.authorize"}\n`
	_ = init2
	_ = init1
	fmt.Println([]byte(init1))
	fmt.Println(hex.EncodeToString([]byte(init1)))
	// fmt.Println(
	writeData(w, []byte(init1[:len(init1)]))
	
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	// rdr := bufio.NewReader(conn)
	fmt.Println("Read...")
	// read, err := r.ReadByte()
	if err != nil {
		fmt.Println("read:", err)
	}
	// _ = buf
	// fmt.Println(read)
	fmt.Println(buf)
	// if read != 0 {
	// 	fmt.Println("Data:", string(buf))
	// } else {
	// 	fmt.Println("No data read")
	// }

	// err = json.Unmarshal()
}

func writeData(w *bufio.Writer, data []byte) {
	w.Write(data)
	w.Flush()
}

func readData(r *bufio.Reader) []byte {
	x := make([]byte, 0)

	return x
}


/*
def get_current_block_height():
	# returns the current network height
	r = requests.get('https://blockchain.info/latestblock')
	return int(r.json()['height'])


def bitcoin_miner(t , restarted = False):

	target = (ctx.nbits[2:] + '00' * (int(ctx.nbits[:2], 16) - 3)).zfill(64)
	randInt = random.randint(0, 2**32-1)
	extranonce2 = hex(randInt)[2:].zfill(2 * ctx.extranonce2_size)

	coinbase = ctx.coinb1 + ctx.extranonce1 + extranonce2 + ctx.coinb2
	coinbase_hash_bin = hashlib.sha256(hashlib.sha256(binascii.unhexlify(coinbase)).digest()).digest()
	merkle_root = coinbase_hash_bin
	for h in ctx.merkle_branch:
		merkle_root = hashlib.sha256(hashlib.sha256(merkle_root + binascii.unhexlify(h)).digest()).digest()
	merkle_root = binascii.hexlify(merkle_root).decode()
	merkle_root = ''.join([merkle_root[i] + merkle_root[i + 1] for i in range(0 , len(merkle_root) , 2)][: :-1])
	work_on = get_current_block_height()
	ctx.nHeightDiff[work_on + 1] = 0
	_diff = int("00000000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF" , 16)

	print(Fore.MAGENTA, '[', timer(), ']', Fore.YELLOW, '[*] Working to solve block with ', Fore.RED, 'height {}'.format(work_on + 1))
	while True:
		t.check_self_shutdown()
		if t.exit:
			break
		if ctx.prevhash != ctx.updatedPrevHash:
			logg('[*] New block {} detected on network '.format(ctx.prevhash))
			print(Fore.YELLOW, '[', timer(), ']', Fore.MAGENTA, '[*] New block {} detected on', Fore.BLUE, ' network '.format(ctx.prevhash))
			logg('[*] Best difficulty will trying to solve block {} was {}'.format(work_on + 1 , ctx.nHeightDiff[work_on + 1]))
			print(Fore.MAGENTA, '[', timer(), ']', Fore.GREEN, '[*] Best difficulty will trying to solve block', Fore.WHITE, '{}', Fore.BLUE, 'was {}'.format(work_on + 1, ctx.nHeightDiff[work_on + 1]))
			ctx.updatedPrevHash = ctx.prevhash
			bitcoin_miner(t , restarted = True)
			print(Back.YELLOW , Fore.MAGENTA , '[' , timer() , ']' , Fore.BLUE , 'Bitcoin Miner Restart Now...' , Style.RESET_ALL)
			continue

		nonce = hex(random.randint(0 , 2 ** 32 - 1))[2 :].zfill(8)  # nNonce   #hex(int(nonce,16)+1)[2:]
		blockheader = ctx.version + ctx.prevhash + merkle_root + ctx.ntime + ctx.nbits + nonce + '000000800000000000000000000000000000000000000000000000000000000000000000000000000000000080020000'
		hash = hashlib.sha256(hashlib.sha256(binascii.unhexlify(blockheader)).digest()).digest()

		hash = binascii.hexlify(hash).decode()
		if hash.startswith('0000000'):
			logg('[*] New hash: {} for block {}'.format(hash , work_on + 1))
			# print(Fore.MAGENTA, '[', timer() , ']', Fore.YELLOW, '[*] New hash:', Fore.WHITE, '{} for block', Fore.WHITE, {}'.format(hash , work_on + 1)
			print(Fore.MAGENTA, '[', timer(), ']', Fore.BLUE, 'Hash:', str(hash))
		this_hash = int(hash , 16)
		difficulty = _diff / this_hash
		if ctx.nHeightDiff[work_on + 1] < difficulty:
			ctx.nHeightDiff[work_on + 1] = difficulty
		if hash < target:
			logg('[*] Block {} solved.'.format(work_on + 1))
			print(Fore.MAGENTA , '[' , timer() , ']' , Fore.YELLOW , '[*] Block {} solved.'.format(work_on + 1))
			logg('[*] Block hash: {}'.format(hash))
			print(Fore.YELLOW)
			print(Fore.MAGENTA , '[' , timer() , ']' , Fore.YELLOW , '[*] Block hash: {}'.format(hash))
			logg('[*] Blockheader: {}'.format(blockheader))
			print(Fore.YELLOW , '[*] Blockheader: {}'.format(blockheader))
			payload = bytes('{"params": ["'+address+'","'+ctx.job_id+'","'+ctx.extranonce2+'","'+ctx.ntime+'","'+nonce+'"],"id":1,"method":"mining.submit"}\n','utf-8')
			logg('[*] Payload:{}'.format(payload))
			print(Fore.MAGENTA, '[' , timer() , ']' , Fore.BLUE , '[*] Payload:' , Fore.GREEN , ' {}'.format(payload))
			sock.sendall(payload)
			ret = sock.recv(1024)
			logg('[*] Pool response: {}'.format(ret))
			print(Fore.MAGENTA, '[', timer(), ']', Fore.GREEN, '[*] Pool Response:', Fore.CYAN, '{}'.format(ret))
			return True


def block_listener(t):
	lines = sock.recv(1024).decode().split('\n')
	response = json.loads(lines[0])
	print(response)
	ctx.sub_details, ctx.extranonce1, ctx.extranonce2_size = response['result']
	print(ctx.sub_details, ctx.extranonce1, ctx.extranonce2_size)
	sock.sendall(b'{"params": ["' + address.encode() + b'", "password"], "id": 2, "method": "mining.authorize"}\n')
	response = b''
	while response.count(b'\n') < 4 and not (b'mining.notify' in response):
		response += sock.recv(1024)
		print(response)
	responses = [json.loads(res) for res in response.decode().split('\n') if len(res.strip()) > 0 and 'mining.notify' in res]
	print(responses)
	ctx.job_id , ctx.prevhash , ctx.coinb1 , ctx.coinb2 , ctx.merkle_branch , ctx.version , ctx.nbits , ctx.ntime , ctx.clean_jobs = responses[0]['params']
	ctx.updatedPrevHash = ctx.prevhash
	while True:
		t.check_self_shutdown()
		if t.exit:
			break
		response = b''
		while response.count(b'\n') < 4 and not (b'mining.notify' in response) : response += sock.recv(1024)
		responses = [json.loads(res) for res in response.decode().split('\n') if len(res.strip()) > 0 and 'mining.notify' in res]
		if responses[0]['params'][1] != ctx.prevhash:
			ctx.job_id , ctx.prevhash , ctx.coinb1 , ctx.coinb2 , ctx.merkle_branch , ctx.version , ctx.nbits , ctx.ntime , ctx.clean_jobs = responses[0]['params']



sock.sendall(b'{"id": 1, "method": "mining.subscribe", "params": []}\n')


{'result': [[['mining.notify', '693da59c']], '1a5a3d69', 8], 'id': 1, 'error': None}
[['mining.notify', '693da59c']] 1a5a3d69 8
b'{"result":false,"error":"Empty workername parameter","id":2}\n'

sock.sendall(b'{"params": ["' + address.encode() + b'", "password"], "id": 2, "method": "mining.authorize"}\n')





*/
