# Existing Issues and Challenges

This section documents the current issues identified in the EasyPars project during the parsing of fight data from https://vringe.com/results/. These issues are logged as of July 05, 2025, and will be addressed in future iterations.

## Identified Issues
1. **Date Parsing Errors**
   - **Description**: Errors like `strconv.Atoi: parsing "": invalid syntax` occur when empty `<td class="date">` values are encountered.
   - **Impact**: Some fight records may be skipped or incorrectly dated.
   - **Example**: Log entry "Error parsing day number '': strconv.Atoi: parsing "": invalid syntax".

2. **Duplicate IDs**
   - **Description**: All fight IDs use the same timestamp (e.g., "fight_1751700163664602800_X") with an incremental counter.
   - **Impact**: Potential ID collisions if parsing multiple pages or months concurrently.
   - **Example**: JSON output shows "id": "fight_1751700163664602800_1" for multiple records.

3. **Missing Title Data**
   - **Description**: Additional `<tr>` rows with title information (e.g., "Бой за титул...") are not extracted.
   - **Impact**: Loss of contextual data about fight significance.
   - **Example**: HTML contains "Бой за титул..." but JSON lacks a "title" field.

4. **Unknown Locations**
   - **Description**: Some fights are assigned "Unknown Location" due to missing `<td class="place">`.
   - **Impact**: Incomplete geographic context for certain events.
   - **Example**: JSON shows "location": "Unknown Location" for some records.

## Causes
- **Date Parsing Errors**: Lack of validation for empty or invalid `<td class="date">` values in `extractFightData`.
- **Duplicate IDs**: Single timestamp generation at the start of parsing, not per fight or location.
- **Missing Title Data**: `extractFightElements` and `extractFightData` do not process additional `<tr>` rows beyond main fight data.
- **Unknown Locations**: No fallback mechanism for missing `<td class="place">` data.

## Potential Solutions
- **Date Parsing Errors**: Add validation to handle empty dates with a default (e.g., current date) or skip with detailed logging.
- **Duplicate IDs**: Implement per-fight or per-location timestamp + counter for unique IDs.
- **Missing Title Data**: Extend parsing logic to extract text from additional `<tr>` rows into a "title" field.
- **Unknown Locations**: Set "N/A" for missing locations and log the issue for review.

## Next Steps
- Document additional edge cases as they are identified.
- Plan resolution in the next development cycle after GitHub commit.
- Update this section with resolution status post-implementation.