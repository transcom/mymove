// If any required fields on the PPM shipment are missing, return as incomplete (true)
export default function checkPPMCompletion(ppm) {
  if (ppm.sitExpected === undefined || ppm.sitExpected === null) return true;
  if (!ppm.expectedDepartureDate) return true;
  if (!ppm.pickupPostalCode) return true;
  if (!ppm.destinationPostalCode) return true;
  return false;
}
