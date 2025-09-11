import os
import requests
import zipfile
import shutil
from tqdm import tqdm
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.chrome.service import Service
from . import inject_captcha_solver

def download_chrome_and_driver(driver_path, version="140.0.7339.82"):
    chrome_url = f"https://storage.googleapis.com/chrome-for-testing-public/{version}/linux64/chrome-linux64.zip"
    driver_url = f"https://storage.googleapis.com/chrome-for-testing-public/{version}/linux64/chromedriver-linux64.zip"
    
    chrome_dir = os.path.join(driver_path, "chrome-linux64")
    driver_dir = os.path.join(driver_path, "chromedriver-linux64")
    
    # Download and extract Chrome if not exists
    if not os.path.exists(os.path.join(chrome_dir, "chrome")):
        chrome_zip_path = os.path.join(driver_path, "chrome-linux64.zip")
        print(f"Downloading Chrome for Testing version {version}...")
        with requests.get(chrome_url, stream=True) as r:
            total_size = int(r.headers.get('content-length', 0))
            with open(chrome_zip_path, 'wb') as f:
                with tqdm(total=total_size, unit='B', unit_scale=True, desc="Downloading Chrome") as pbar:
                    for chunk in r.iter_content(chunk_size=1024):
                        if chunk:
                            f.write(chunk)
                            pbar.update(len(chunk))
        print("Extracting Chrome...")
        with zipfile.ZipFile(chrome_zip_path, 'r') as zip_ref:
            zip_ref.extractall(driver_path)
        os.remove(chrome_zip_path)
        print("Chrome downloaded and extracted.")
    
    # Download and extract Chromedriver if not exists
    if not os.path.exists(os.path.join(driver_dir, "chromedriver")):
        driver_zip_path = os.path.join(driver_path, "chromedriver-linux64.zip")
        print(f"Downloading Chromedriver version {version}...")
        with requests.get(driver_url, stream=True) as r:
            total_size = int(r.headers.get('content-length', 0))
            with open(driver_zip_path, 'wb') as f:
                with tqdm(total=total_size, unit='B', unit_scale=True, desc="Downloading Chromedriver") as pbar:
                    for chunk in r.iter_content(chunk_size=1024):
                        if chunk:
                            f.write(chunk)
                            pbar.update(len(chunk))
        print("Extracting Chromedriver...")
        with zipfile.ZipFile(driver_zip_path, 'r') as zip_ref:
            zip_ref.extractall(driver_path)
        os.remove(driver_zip_path)
        print("Chromedriver downloaded and extracted.")

def build_config(driver_path="driver", version="140.0.7339.82"):
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
    
    options = inject_captcha_solver(options)
    service = Service(chromedriver_binary_path)
    return options, service