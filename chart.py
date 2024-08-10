import json
import matplotlib.pyplot as plt

data = '''[{}]''' # Insert the JSON data here

# Load the JSON data
json_data = json.loads(data)

# Extract data for plotting
days = [entry["Day"] for entry in json_data]
attended_patients = [entry["AttendedPacients"] for entry in json_data]
not_attended_patients = [entry["NotAttendedPacients"] for entry in json_data]

# Plot the data
plt.plot(days, attended_patients, label='Attended Patients')
plt.plot(days, not_attended_patients, label='Not Attended Patients')
plt.xlabel('Days')
plt.ylabel('Number of Patients')
plt.title('COVID Simulation Results')
plt.legend()
plt.show()
