import csv
import glob
stats = {
  "Name": [],
  "AverageSize": [],
  "MedianSize": [],
  "StdDev": [],
  "TotalSize": [],
  "TotalObjects": []
}


total_files = 0
for filename in glob.glob("*.oracleGeneral"):
  objects = {}
  total_bytes=0
  total_unique_objects=0
  with open(filename, "rb") as f:
    while True:
      line = f.read(24)
      if len(line) < 24:
         break
      
      objId = int.from_bytes(line[4:12], "little")
      objSize = int.from_bytes(line[12:16], "little")
          
      if objId not in objects:
        total_unique_objects += 1
        total_bytes += objSize
        objects[objId] = objSize
  f.close()
  total_files += 1
  average = total_bytes/total_unique_objects   
  stats["Name"].append(filename.split(".")[0])
  stats["AverageSize"].append(average)
  stats["TotalObjects"].append(total_unique_objects)
  stats["TotalSize"].append(total_bytes)

with open("stats.csv", "w", newline="") as f:
    writer = csv.writer(f)
    # header
    writer.writerow(stats.keys())
    for i in range(total_files):
        row = []
        for key in stats:
          print(key,i,len(stats[key]))
          row.append(stats[key][i] if len(stats[key]) > 0 else "")
        writer.writerow(row)