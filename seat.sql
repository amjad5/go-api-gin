SELECT 
    CASE 
        WHEN ID % 2 = 1 AND ID < (SELECT MAX(ID) FROM Seat) THEN ID + 1
        WHEN ID % 2 = 0 THEN ID - 1
        ELSE ID
    END AS ID,
    student
FROM Seat
ORDER BY ID;