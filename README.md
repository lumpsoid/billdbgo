# Bill Manager

## Overview

This project is a web application and API built in Go (Golang) for managing bills by parsing QR codes from Serbian bills. It allows users to easily scan their bills, extract relevant information, and manage their expenses efficiently.

You can use flutter app to scan QR codes and send them to the server. The flutter app can be found here: [qr_bills](https://github.com/lumpsoid/qr_bills).

This is a refactor of the python version of the project. The original project can be found here: [billdb flask server](https://github.com/lumpsoid/billdb_flask_api), and [billdb package](https://github.com/lumpsoid/billdb).

## Features
- Parsing of invoice information: The application can extract information from the QR code on the bill.
- Server-side decoding of QR codes
- Invoice management
    - view
    - search
    - organize
- Expense tracking: Users can track their expenses over time.
- API endpoints: Provides API endpoints for integrating with other applications or services.
