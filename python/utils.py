import re
from typing import List, Callable, Any

DIGIT_REGEXP = re.compile('([0-9]+)')


def natural_sort_key(key: str, _nsre=DIGIT_REGEXP) -> List[str]:
    """
    natural_sort_key can be used as a key in any function that uses it, like `list.sort`, `sorted`, `max`, etc
    :param key: list of string for sorting
    :param _nsre: regular expression for split key
    :type  key: str
    :type  _nsre: regexp
    :return: list of
    :rtype: str
    """
    return [int(text) if text.isdigit() else text.lower()
            for text in re.split(_nsre, key)]


def sorted_natural_order(data: List[Any], key: Callable=lambda s: s) -> List[Any]:
    """
    Sort the list into natural alphanumeric order..
    Return sorted list
    :param data: list of string for sorting
    :param key: key function like for `sorted`
    :type  data: list[str]
    :type  key: callable
    :return: sorted list
    :rtype: list[str]
    """
    return sorted(data, key=lambda x: natural_sort_key(key(x)))


def divide_chunks(l: List[Any], n):
    for i in range(0, len(l), n):
        yield l[i:i + n]