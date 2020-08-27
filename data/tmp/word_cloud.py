from wordcloud import WordCloud, ImageColorGenerator
import jieba
import matplotlib.pyplot as plt
import base64
import sys

filename = sys.argv[1]
savefile=sys.argv[2]

with open(filename, encoding='UTF-8') as f1:
	data1 = f1.read()

wordList_jieba1 = jieba.cut(data1, cut_all=False)

data1 = ','.join(wordList_jieba1)

wc1 = WordCloud(font_path="data/tmp/STKAITI.TTF").generate(data1)

wc1.to_image().save(savefile)
