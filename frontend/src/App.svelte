<script>
	import "../node_modules/@picocss/pico/css/pico.min.css";
	import { SvelteToast, toast } from '@zerodevx/svelte-toast';

	const ACTIONS_PATH = `API_BASE_PATH/whitepapers`;
	const UPDATE_URL = `API_BASE_PATH/whitepaper`;
	const UNREAD_DATE = "0001-01-01T00:00:00Z";

	// Sourcedata is the full list. filteredactions will hold the current results based on the filter
	let sourceData = [];
	let filteredData = [];
	let paperTypes = [];

	let filters = {
		Type: "all",
		State: "all"
	};

	// Load the source iam list
	async function fetchSourceData() {
		const res = await fetch(ACTIONS_PATH);
		sourceData = await res.json();

		sourceData.forEach(wp => {
			if (!paperTypes.includes(wp.Type)) {
				paperTypes = [...paperTypes, wp.Type];
			}
		});
		filteredData = sourceData;
	}

	function filterData() {
		console.log("filtering");
		console.log(filters);
		let newFilteredData = sourceData;
		if (filters.Type != "all") {
			newFilteredData = newFilteredData.filter(p => p.Type == filters.Type);
		}
		if (filters.State != "all") {
			newFilteredData = newFilteredData.filter(p => {
				if (filters.State == "updated") {
					return updatedSinceRead(p);
				} else if (filters.State == "read") {
					return isRead(p.DateRead);
				} else if (filters.State == "unread") {
					return !isRead(p.DateRead);
				}
				console.log("Unknown filter for state, returning all");
				return true;
			});
		}
		filteredData = newFilteredData;
	}

	function isRead(d) {
		return d != UNREAD_DATE;
	}

	function updatedSinceRead(entry) {
		// We haven't read this yet, so there are no updates
		if (!isRead(entry.DateRead)) {
			return false;
		}
		const readDate = Date.parse(entry.DateRead);
		const updateDate = Date.parse(entry.DateUpdated);
		return updateDate > readDate;
	}

	async function markPaperState(id, readState) {
		const result = await fetch(UPDATE_URL, {
			method: 'POST',
			mode: 'cors',
			cache: 'no-cache',
			headers: {
				'Content-Type': 'application/json'
			},
			redirect: "follow",
			body: JSON.stringify({
				Id: id,
				Read: readState
			})
		});
		const newItem = await result.json();
		sourceData[sourceData.findIndex(p => p.Id == newItem.Id)] = newItem;
		filteredData[filteredData.findIndex(p => p.Id == newItem.Id)] = newItem;
	}

	fetchSourceData();
</script>

<SvelteToast />
<main class="container">
	<div>
		<h2>AWS Whitepaper read tracker</h2>
		<p>
			This page tracks which aws papers you've read, and whether any have been updated since you read them
		</p>
	</div>
	<div>
		<form>
			<div class="grid">
				<label for="type">
					Type
					<select id="type" bind:value={filters.Type} on:change="{filterData}">
						<option value="all">All</option>
						{#each paperTypes as pt}
							<option value="{pt}">{pt}</option>
						{/each}
					</select>
				</label>
				<label for="state">
					State
					<select id="state" bind:value={filters.State} on:change="{filterData}">
						<option value="all">All</option>
						<option value="read">Read</option>
						<option value="unread">Unread</option>
						<option value="updated">Updated since last read</option>
					</select>
				</label>
		</form>
	</div>
	<table>
		<thead>
			<tr>
				<td>
					Paper
				</td>
				<td>
					Read
				</td>
			</tr>
		</thead>
		<tbody>
			{#each filteredData as wp,idx}
				<tr>
					<td>
						<a href="{wp.Url}" target="_blank">{wp.Title}</a>
						{#if updatedSinceRead(wp)}
							<mark>Updated</mark>
						{/if}
						<div>
							<small>
								<b>Type: </b>{wp.Type} <br />
								<b>Published: </b>{wp.DatePublished.substr(0, 10)} <b>Last update: </b>{ wp.DatePublished == wp.DateUpdated ? "-" : wp.DateUpdated.substr(0, 10)}
								{#if isRead(wp.DateRead)}
									<br /><b>Read:</b> {wp.DatePublished.substr(0, 10)}
								{/if}
							</small>
						</div>
					</td>
					<td>
						<input on:click={() => markPaperState(wp.Id, !isRead(wp.DateRead))} type="checkbox" checked={isRead(wp.DateRead)}>
					</td>
				</tr>
			{/each}
		</tbody>
	</table>
</main>

<style></style>