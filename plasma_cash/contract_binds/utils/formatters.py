def normalize(d):
    if type(d) is str:
        return bytes(d, 'utf-8')
    elif type(d) is list:
        return [bytes(x, 'utf-8') for x in d]
    else:
        return d
