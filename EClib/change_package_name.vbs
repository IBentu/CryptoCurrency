Const ForReading = 1
Const ForWriting = 2
const direc = "eclib.go"
file = ""
Set objFS = CreateObject("Scripting.FileSystemObject")
Set objFile = objFS.OpenTextFile _
    (direc, ForReading)
do until objFile.AtEndOfStream
    strLine = objFile.ReadLine
    If InStr(strLine,"package eclib")> 0 Then
        file = file + "package main" + vbCrLf
    Else
        file = file + strLine + vbCrLf
    End If 
loop
objFile.Close

Set objFile = objFS.OpenTextFile _
    (direc, ForWriting)
objFile.Write(file)
objFile.Close