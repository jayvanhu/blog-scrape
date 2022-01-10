from typing import List
from channel import Channel

# TODO add type to process?
def pipe_thru_channel(input_q: Channel, output_q: Channel, process):
	'''
	Reads data from `input_q`, processes them, then puts them into `output_q`.
	Can be run concurrently.
	'''
	while not input_q.finished():
		data = input_q.get()
		processed = process(data)
		output_q.put(processed)
		input_q.task_done()

def pipe_thru_channel_many(input_q: Channel, output_q: Channel, process):
	'''
	Reads data from `input_q`, processes them into an iterable, then puts items of iterable into `output_q`.
	Can be run concurrently.
	'''
	while not input_q.finished():
		data = input_q.get()
		processed_items = process(data)
		for item in processed_items:
			output_q.put(item)
		input_q.task_done()

def send_to_channel(items: List, chan: Channel):
	'''
	Sends all the contents of `items` into `chan`.
	Can be run concurrently.
	'''
	for item in items:
		chan.put(item)
	print('Close sending chan')
	chan.close()

class ChannelReceiver:
	'''(Concurrently) collects'''

	def recv_from_channel(self, chan: Channel) -> List:
		'''
		Reads data from `queue` and collects them into a list.
		Can be run concurrently.
		'''
		self.result = []
		while not chan.finished():
			print('Recv from channel')
			data = chan.get()
			print('Recv data: ', data)
			self.result.append(data)
			chan.task_done()

	# TODO apply fail safe or check if channel has not been fully processed?
	def list(self) -> List:
		return self.result
