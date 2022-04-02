# CSV Splitter

CSV Splitter is Go Applications, Used to Split CSV Files into smaller chunks.

## No Installation Needed

Just Double click on the Application and you are good to go.

## Features

* Auto Detect CSV file in current folder
* Split by no. of Rows (or) Split by no. of files
* Option to store each chunk file in separate folders

## Detailed Explanation

![image](https://user-images.githubusercontent.com/47551990/161386492-7b061d97-81b8-40b8-a493-91b5cdd4d71e.png)

1. CSV File will be Autodetected and FileName will be displayed at the Top
2. Chunk by have top options (Number of Rows or Number of Chunks)
    -> Select Number of Rows if you like to split CSV file based of rows (All files no. of row you provided)
    -> Select Number of Chunks if you like to split CSV file based on no. of files (All files row will be calculated based on your provided value)
3. Next you should provide value based on the option you chossed earlier
4. Optionally you can select Create folder for each chunk to store each chuck in seperate folder.
5. Click Split Now.

## Note
For all error, you will be notified in detail.

## Example
If you have CSV file with 100 rows, and you like to split chucks each with 10 rows.
1. Selected Chuck by -> Number of Rows 
2. Enter 10 in Enter Value field
3. Click Split Now

You can see 10 new files will be created each with 10 rows.
