from invoke import task


@task
def verify(c):
    c.run("black --check .")
    c.run("flake8")
    c.run("python -m unittest")
