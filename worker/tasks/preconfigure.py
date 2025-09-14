import os
import requests
import zipfile
import shutil
from tqdm import tqdm
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from . import inject_captcha_solver
import subprocess

def download_chrome_and_driver(driver_path, version="140.0.7339.82"):
    os.makedirs(driver_path, exist_ok=True)
    chrome_url = f"https://storage.googleapis.com/chrome-for-testing-public/{version}/linux64/chrome-linux64.zip"
    driver_url = f"https://storage.googleapis.com/chrome-for-testing-public/{version}/linux64/chromedriver-linux64.zip"
    
    chrome_dir = os.path.join(driver_path, "chrome-linux64")
    driver_dir = os.path.join(driver_path, "chromedriver-linux64")
    
    # Download and extract Chrome if not exists
    if not os.path.exists(os.path.join(chrome_dir, "chrome")):
        chrome_zip_path = os.path.join(driver_path, "chrome-linux64.zip")
        print(f"Downloading Chrome for Testing version {version}...")
        subprocess.run(["wget", chrome_url, "-O", chrome_zip_path], check=True)
        print("Extracting Chrome...")
        subprocess.run(["unzip", chrome_zip_path, "-d", driver_path], check=True)
        os.remove(chrome_zip_path)
        print("Chrome downloaded and extracted.")
    
    # Set permissions for Chrome
    chrome_path = os.path.join(chrome_dir, "chrome")
    if os.path.exists(chrome_path):
        os.chmod(chrome_path, 0o755)
    
    # Download and extract Chromedriver if not exists
    if not os.path.exists(os.path.join(driver_dir, "chromedriver")):
        driver_zip_path = os.path.join(driver_path, "chromedriver-linux64.zip")
        print(f"Downloading Chromedriver version {version}...")
        subprocess.run(["wget", driver_url, "-O", driver_zip_path], check=True)
        print("Extracting Chromedriver...")
        subprocess.run(["unzip", driver_zip_path, "-d", driver_path], check=True)
        os.remove(driver_zip_path)
        print("Chromedriver downloaded and extracted.")
    
    # Set permissions for Chromedriver
    driver_path_bin = os.path.join(driver_dir, "chromedriver")
    if os.path.exists(driver_path_bin):
        os.chmod(driver_path_bin, 0o755)

def build_config(driver_path="driver", version="140.0.7339.82", headless=False):
    options = Options()
    current_dir = os.path.dirname(__file__)
    upper_dir = os.path.dirname(current_dir)
    
    if os.path.isabs(driver_path):
        full_driver_path = driver_path
    else:
        full_driver_path = os.path.join(upper_dir, driver_path)
    
    download_chrome_and_driver(full_driver_path, version)
    
    chrome_binary_path = os.path.join(full_driver_path, "chrome-linux64", "chrome")
    chromedriver_binary_path = os.path.join(full_driver_path, "chromedriver-linux64", "chromedriver")
    options.binary_location = chrome_binary_path
    
    # Добавляем опции для имитации обычного пользователя
    options.add_argument("--disable-blink-features=AutomationControlled")
    options.add_experimental_option("excludeSwitches", ["enable-automation"])
    options.add_experimental_option('useAutomationExtension', False)
    options.add_argument("--user-agent=Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
    if headless:
        options.add_argument("--headless=new")
    
    options = inject_captcha_solver(options)
    service = Service(chromedriver_binary_path)
    return options, service