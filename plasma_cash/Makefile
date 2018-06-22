test:
	pytest

lint:
	flake8 child_chain client contract_binds utils *.py --ignore=E203,W503

format:
	black --line-length 79 -S .
	isort -rc .
