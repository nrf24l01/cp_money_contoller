from selenium.webdriver.chrome.options import Options
import os

def inject_captcha_solver(options: Options):
    extension_path = os.path.join(os.path.dirname(__file__), "../captcha_solver.crx")
    options.add_extension(extension_path)
    return options