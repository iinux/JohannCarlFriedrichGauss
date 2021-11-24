import holoviews as hv
from holoviews import opts, dim
from bokeh.sampledata.airport_routes import routes, airports
import time

# 选择bokeh引擎
hv.extension('bokeh')

# Count the routes between Airports
route_counts = routes.groupby(
    ['SourceID', 'DestinationID']).Stops.count().reset_index()
nodes = hv.Dataset(airports, 'AirportID', 'City')
chord = hv.Chord((route_counts, nodes), [
                 'SourceID', 'DestinationID'], ['Stops'])

# Select the 6 busiest airports
busiest = list(routes.groupby('SourceID').count(
).sort_values('Stops').iloc[-6:].index.values)
busiest_airports = chord.select(AirportID=busiest, selection_mode='nodes')

busiest_airports.opts(
    opts.Chord(cmap='Category20', edge_color=dim('SourceID').str(),
               height=500,
               labels='City',
               node_color=dim('AirportID').str(), width=500))

hv.save(busiest_airports, r'output.html')
