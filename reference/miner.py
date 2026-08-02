import binascii
import hashlib
import json
import logging
import random
import socket
import threading
import time
import traceback
import requests
from datetime import datetime
from signal import SIGINT , signal
from colorama import Back , Fore , Style
import contextvars as ctx
# import context as ctx


sock = None

def timer():
	tcx = datetime.now().time()
	return tcx

pool = ('solo.ckpool.org' , 3333)
server = ("127.0.0.1", 5656)
address = 'bc1q03zbvbamcx77crptfga5e88cmubnea0jctzdfq'

print(Back.BLUE , Fore.WHITE , 'WALLET:' , Fore.BLACK , str(address) , Style.RESET_ALL)

def handler(signal_received , frame):
	ctx.fShutdown = True
	print(Fore.MAGENTA , '[' , timer() , ']' , Fore.YELLOW , 'Terminating Miner, Please Wait..')
	quit()

def logg(msg):
	logging.basicConfig(level = logging.INFO , filename = "miner.log", format = '%(asctime)s %(message)s')  # include timestamp
	logging.info(msg)

def get_current_block_height():
	# returns the current network height
	r = requests.get('https://blockchain.info/latestblock')
	return int(r.json()['height'])

def bitcoin_miner(t , restarted = False):
	if restarted:
		logg('\n[*] Bitcoin Miner restarted')
		print(Fore.MAGENTA , '[' , timer() , ']' , Fore.YELLOW , 'Programmer = Mmdrza.Com')
		print(Fore.MAGENTA , '[' , timer() , ']' , Fore.BLUE , '[*] Bitcoin Miner Restarted')
		time.sleep(5)
	
	target = (ctx.nbits[2:] + '00' * (int(ctx.nbits[:2], 16) - 3)).zfill(64)
	randInt = random.randint(0, 2**32-1)
	extranonce2 = hex(randInt)[2:].zfill(2 * ctx.extranonce2_size)
	print(extranonce2)
	
	coinbase = ctx.coinb1 + ctx.extranonce1 + extranonce2 + ctx.coinb2
	coinbase_hash_bin = hashlib.sha256(hashlib.sha256(binascii.unhexlify(coinbase)).digest()).digest()
	merkle_root = coinbase_hash_bin
	for h in ctx.merkle_branch:
		merkle_root = hashlib.sha256(hashlib.sha256(merkle_root + binascii.unhexlify(h)).digest()).digest()
	merkle_root = binascii.hexlify(merkle_root).decode()
	merkle_root = ''.join([merkle_root[i] + merkle_root[i + 1] for i in range(0 , len(merkle_root) , 2)][: :-1])
	work_on = get_current_block_height()
	ctx.nHeightDiff[work_on + 1] = 0
	#000000000000000000013f3344a94d3383c4b97f92d51dd283514a58a035b21f
	_diff = int("0000000000000000000FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF" , 16)
	
	logg('[*] Working to solve block with height {}'.format(work_on + 1))
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
		print(blockheader)
		send_data(blockheader)
		hash = hashlib.sha256(hashlib.sha256(binascii.unhexlify(blockheader)).digest()).digest()
	
		hash = binascii.hexlify(hash).decode()
		if hash.startswith('0000000000000000000'):
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
	sock = socket.socket(socket.AF_INET , socket.SOCK_STREAM)
	sock.connect(pool)
	sock.sendall(b'{"id": 1, "method": "mining.subscribe", "params": []}\n')
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

def StartMining():
	subscribe_t = NewSubscribeThread(None)
	subscribe_t.start()
	logg("[*] Subscribe thread started.")
	print(Fore.MAGENTA , "[" , timer() , "]" , Fore.GREEN , "[*] Subscribe thread started.")
	time.sleep(4)
	miner_t = CoinMinerThread(None)
	miner_t.start()
	logg("[*] Bitcoin Miner Thread Started")
	print(Fore.MAGENTA, "[", timer(), "]", Fore.GREEN, "[*] Thread Started")
	print(Fore.BLUE, '--------------~~(' , Fore.YELLOW , '' , Fore.BLUE , ' )~~--------------')

def check_for_shutdown(t):
	n = t.n
	if ctx.fShutdown:
		if n != -1 :
			ctx.listfThreadRunning[n] = False
			t.exit = True

def send_data(d):
	sock = socket.socket(socket.AF_INET , socket.SOCK_STREAM)
	sock.connect(server)
	sock.sendall(d)

class ExitedThread(threading.Thread):
	def __init__(self , arg , n) :
		super(ExitedThread , self).__init__()
		self.exit = False
		self.arg = arg
		self.n = n

	def run(self):
		self.thread_handler(self.arg , self.n)
		pass

	def thread_handler(self , arg , n):
		while True:
			check_for_shutdown(self)
			if self.exit:
				break
			ctx.listfThreadRunning[n] = True
			try:
				self.thread_handler2(arg)
			except Exception as e:
				logg("ThreadHandler()")
				print(Fore.MAGENTA , '[' , timer() , ']' , Fore.WHITE , 'ThreadHandler()')
				logg(e)
				print(Fore.RED , e)
			ctx.listfThreadRunning[n] = False
			time.sleep(2)
			pass

	def thread_handler2(self , arg):
		raise NotImplementedError("must impl this func")

	def check_self_shutdown(self):
		check_for_shutdown(self)

	def try_exit(self):
		self.exit = True
		ctx.listfThreadRunning[self.n] = False
		pass

class CoinMinerThread(ExitedThread):
	def __init__(self , arg = None):
		super(CoinMinerThread , self).__init__(arg , n = 0)

	def thread_handler2(self , arg):
		self.thread_bitcoin_miner(arg)

	def thread_bitcoin_miner(self , arg):
		ctx.listfThreadRunning[self.n] = True
		check_for_shutdown(self)
		try:
			ret = bitcoin_miner(self)
			logg(Fore.MAGENTA , "[" , timer() , "] [*] Miner returned %s\n\n" % "true" if ret else "false")
			print(Fore.LIGHTCYAN_EX , "[*] Miner returned %s\n\n" % "true" if ret else "false")
		except Exception as e:
			logg("[*] Miner()")
			print(Back.WHITE , Fore.MAGENTA , "[" , timer() , "]" , Fore.BLUE , "[*] Miner()")
			logg(e)
			traceback.print_exc()
		ctx.listfThreadRunning[self.n] = False

class NewSubscribeThread(ExitedThread):
	def __init__(self , arg = None):
		super(NewSubscribeThread , self).__init__(arg , n = 1)

	def thread_handler2(self , arg):
		self.thread_new_block(arg)

	def thread_new_block(self , arg):
		ctx.listfThreadRunning[self.n] = True
		check_for_shutdown(self)
		try:
			ret = block_listener(self)
		except Exception as e:
			logg("[*] Subscribe thread()")
			print(Fore.MAGENTA , "[" , timer() , "]" , Fore.YELLOW , "[*] Subscribe thread()")
			logg(e)
			traceback.print_exc()
		ctx.listfThreadRunning[self.n] = False

if __name__ == '__main__':
	signal(SIGINT , handler)
	block_listener(None)
	StartMining()
