---

title: Workflow für Reminder
date: 2021-06-28

---

* Dateien mit der Extension *.rem enthalten Reminder
* eingelesen werden der Reihe nach alle rem-Dateien in einem Verzeichnis

* Lexer geht Zeichen für Zeichen über eine Datei
* ein Zeilenumbruch schließt einen Reminder ab / beginnt einen neuen Reminder
* Lexer bildet Zeichengruppen
* Zeichengruppen können Zahlen, einzelne Worte oder Gruppen von Zahlen und Worten sein
* Blanks, Tabs und Zeilenumbrüche sind Trennzeichen

* Parser gibt den gebildeten Zeichengruppen eine Bedeutung
* steuernde Elemente sind ...
  * Tag, Monat und Jahr oder deren teilweise Abwesenheit
  * einzelne oder mehrere Wochentage
* modifiziert werden die steuernden Elemente durch ...  
  * Wiederholungen alle n Tage
  * vorzeitige Anzeige um n Tage

* Evaluator validiert einzelne steuernde Elemente oder ihre Kombinationen ...
  * Tag, Monat und Jahr bei vollständiger Angabe als Datum
  * Tag und Jahr jeweils bei einzelner Angabe
  * Beginn als Uhrzeit

* Evaluator steuert die Anzeige eines Reminders
* benötigt dazu Beginn und Ende des anzuzeigenden Datumsbereiches
* Ereignisse können innerhalb des Datumsbereiches mehrfach auftreten durch ...
  * Wiederholungen
  * vorzeitige Anzeige
  * Kombination aus beidem
* durch vorzeitige Anzeige kann außerhalb des Anzeigebereiches liegender Reminder sichtbar werden

* Evaluator steuert die farbige Anzeige eines Reminders
* vorderer Teil des Dateinamens kann in Konfiguration mit Farbe belegt werden

* Lexer => Parser => Evaluator

* Layout der Anzeige

12.09.2021
- Tag des Bleistifts
- Geburtstag Gustav
- 09:00 - 09:30 Teamrunde
- 10:00         Workshop zum Reminder-Projekt

13.09.2021
- 09:00 - 09:45 große Teamrunde
- 10:00         Anforderungsanalyse
- 12:00 - 13:30 Mittagessen mit Team
- 19:30 - 21:30 Konzert
