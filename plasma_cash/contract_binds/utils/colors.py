class Colors:
    GREEN = '\033[92m'
    YELLOW = '\033[93m'
    RED = '\033[91m'
    END = '\033[0m'


def green(txt):
    return Colors.GREEN + txt + Colors.END


def yellow(txt):
    return Colors.YELLOW + txt + Colors.END


def red(txt):
    return Colors.RED + txt + Colors.END
