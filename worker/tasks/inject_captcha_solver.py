from selenium.webdriver.chrome.options import Options
import os

def inject_captcha_solver(options: Options):
    extension_path = os.path.join(os.path.dirname(__file__), "../booster")
    options.add_argument("--load-extension=" + extension_path)
    print(extension_path)
    return options