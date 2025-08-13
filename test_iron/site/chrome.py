import webview
import tkinter as tk
from tkinter import ttk

def open_internal_browser():
    webview.create_window("My Local App", "http://127.0.0.1:5000", width=1200, height=800)
    webview.start()

# GUI setup
root = tk.Tk()
root.title("Local Web App Launcher")

btn = ttk.Button(root, text="🧭 เปิดแอพ", command=open_internal_browser)
btn.pack(pady=20)

root.mainloop()
