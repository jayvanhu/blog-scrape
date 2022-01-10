from queue import Queue

class Channel(Queue):
	'''
	A subclass of Queue that adds `close()` and `finished()` methods to get it to behave similar to Golang channels.
	'''
	def __init__(self, maxsize: int = 0) -> None:
		super().__init__(maxsize=maxsize)
		self._closed = False

	# TODO should I throw error if put() is called after close()?
	def close(self) -> None:
		'''Marks a Channel as closed'''
		self._closed = True

	def finished(self) -> bool:
		'''Returns True if a channel is closed and empty'''
		return self._closed and self.empty()
