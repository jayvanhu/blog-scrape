import os
from time import sleep
from bs4 import BeautifulSoup
import json
from threading import Thread
from typing import List
import requests

from article import Article
from channel import Channel
from concurrency import ChannelReceiver, pipe_thru_channel_many, send_to_channel
from config import config

def fetch_archive_links(site_url: str, selector) -> List[str]:
	'''
	Fetches `site_url` page and collects archive links using `selector`.
	An "archive link" links to a list of that archive's articles for the month.
	'''
	# TODO set header?
	# req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:10.0) Gecko/20100101 Firefox/10.0")
	res = requests.get(site_url)
	doc = BeautifulSoup(res.text, 'html.parser')
	anchor_tags = doc.select(selector)
	print('Anchor tags found: ', len(anchor_tags))
	hrefs = map(lambda anchor : anchor['href'], anchor_tags)
	return hrefs

def scrape_articles_from_archive(archive_url: str):
	'''
	Fetches `archive_url` page, collects all the articles on the page into a list, and returns the list.
	Will be the `process` arg for `pipe_thru*` functions
	Intended to be run in multiple threads concurrently, so there needs to be a throttle on network requests.
	'''
	print('Fetch archive_url')
	res = requests.get(archive_url)
	doc = BeautifulSoup(res.text, 'html.parser')

	articles = []
	print('Process archived articles')
	for header in doc.select('header.entry-header'):
		anchor = header.select_one('h1.entry-title > a')
		title = anchor.text
		href = anchor['href']
		date_str = header.select_one('time.entry-date').get('datetime')
		if title and href and date_str:
			articles.append( Article(title, href, date_str) )
		else:
			print(f'{scrape_articles_from_archive.__name__}() :: field missing: {title}, {href}, {date_str}')
	print('Throttle')
	sleep(config.requestDelay)
	return articles

joel_blog_url = 'https://www.joelonsoftware.com/archives/'
# TODO this exact, commented out links_selector from the golang version works perfectly and is valid in js DOM, but only gets the first a tag here?
# links_selector = '.yearly-archive > .month > h3 > a'
links_selector = '.yearly-archive > .month h3 > a'

# TODO unify url vs link variable names
def start_scrape():
	print('Fetch archive links')
	archive_links = fetch_archive_links(joel_blog_url, links_selector)

	if config.fastDebug:
		archive_links = [*archive_links][:2]
		print('Links: ', archive_links)

	## Channels
	urls_q = Channel()
	'''Receives archive links to be scraped for its articles'''
	articles_q = Channel()
	'''Receives finalized Article objects to be collected into a list'''

	## Send
	print('Send to queue')
	Thread(
		target=send_to_channel,
		args=(archive_links, urls_q),
		daemon=True
	).start()

	## Process
	print('Create scraper threads')
	for _i in range(config.threadCount):
		print('Thread: ', _i)
		Thread(target=pipe_thru_channel_many, args=(urls_q, articles_q, scrape_articles_from_archive), daemon=True).start()

	## Receive
	print('Receive articles')
	receiver = ChannelReceiver()
	Thread(
		target=receiver.recv_from_channel,
		args=(articles_q,),
		daemon=True
	).start()

	## Serialize
	print('Wait urls_q')
	urls_q.join()
	print('Close articles_q')
	articles_q.close()
	print('Wait articles_q')
	articles_q.join()
	articles: List[Article] = receiver.list()

	articles_json = [ article.__dict__ for article in articles ]
	file_name = 'dist/scraped-links.json'
	os.makedirs(os.path.dirname(file_name), exist_ok=True)
	with open(file_name, 'w') as file:
		json.dump(articles_json, file, indent=4)

# TODO dockerize app
# TODO generate html representation of scraped data
if __name__ == '__main__':
	start_scrape()
