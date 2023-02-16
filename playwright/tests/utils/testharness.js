// @ts-check
/**
 * Test Harness
 */
export class TestHarness {
  /**
   * Create a TestHarness.
   * @param {import('@playwright/test').APIRequestContext} request
   */
  constructor(request) {
    this.request = request;
  }

  /**
   * @typedef {Object} User
   * @property {string} id
   * @property {string} login_gov_email
   */

  /**
   * @typedef {Object} Move
   * @property {string} id
   * @property {string} locator
   * @property {Object} Orders
   * @property {Object} Orders.NewDutyLocation
   * @property {string} Orders.NewDutyLocation.name
   * @property {Object} Orders.ServiceMember
   * @property {string} Orders.ServiceMember.edipi
   * @property {string} Orders.ServiceMember.last_name
   * @property {string} Orders.ServiceMember.user_id
   * @property {Object} CloseoutOffice
   * @property {string} CloseoutOffice.name
   */

  /**
   * call the test harness build
   *
   * @param {string} action
   * @returns {Promise<Object>} Promise object for created data
   */
  async buildDefault(action) {
    const r = await this.request.post(`/testharness/build/${action}`);
    if (!r.ok()) {
      const body = await r.body();
      throw Error(`Error with testharness build for '${action}': ${body}`);
    }
    const obj = /** @type {Object} */ await r.json();
    await r.dispose();
    return obj;
  }

  /**
   * @returns {Promise<User>}
   */
  async buildDefaultAdminUser() {
    return this.buildDefault('DefaultAdminUser');
  }

  /**
   * build office user with TOO and TIO roles
   * @returns {Promise<User>}
   */
  async buildOfficeUserWithTOOAndTIO() {
    return this.buildDefault('OfficeUserWithTOOAndTIO');
  }

  /**
   * @returns {Promise<User>}
   */
  async buildNeedsOrdersUser() {
    return this.buildDefault('NeedsOrdersUser');
  }

  /**
   * @returns {Promise<Move>}
   */
  async buildDefaultMove() {
    return this.buildDefault('DefaultMove');
  }

  /**
   * @returns {Promise<Move>}
   */
  async buildMoveWithOrders() {
    return this.buildDefault('MoveWithOrders');
  }

  /**
   * @returns {Promise<Move>}
   */
  async buildSpouseProGearMove() {
    return this.buildDefault('SpouseProGearMove');
  }

  /**
   * @returns {Promise<Move>}
   */
  async buildPPMInProgressMove() {
    return this.buildDefault('PPMInProgressMove');
  }

  /**
   * @returns {Promise<Move>}
   */
  async buildMoveWithPPMShipmentReadyForFinalCloseout() {
    return this.buildDefault('MoveWithPPMShipmentReadyForFinalCloseout');
  }

  /**
   * @returns {Promise<Move>}
   */
  async buildPPMMoveWithCloseout() {
    return this.buildDefault('PPMMoveWithCloseout');
  }

  /**
   * @returns {Promise<Move>}
   */
  async buildPPMMoveWithCloseoutOffice() {
    return this.buildDefault('PPMMoveWithCloseoutOffice');
  }

  /**
   * @returns {Promise<Move>}
   */
  async buildApprovedMoveWithPPM() {
    return this.buildDefault('ApprovedMoveWithPPM');
  }

  /**
   * @returns {Promise<Object>}
   */
  async buildWithShipmentMove() {
    return this.buildDefault('WithShipmentMove');
  }

  /**
   * Use testharness to build hhg move for TOO
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO() {
    return this.buildDefault('HHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO');
  }

  /**
   * Use testharness to build retiree hhg move for TOO
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithRetireeForTOO() {
    return this.buildDefault('HHGMoveWithRetireeForTOO');
  }

  /**
   * Use testharness to build HHG move with NTS for TOO
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithNTSShipmentsForTOO() {
    return this.buildDefault('HHGMoveWithNTSShipmentsForTOO');
  }

  /**
   * Use testharness to build move with NTS for TOO
   * @returns {Promise<Move>}
   */
  async buildMoveWithNTSShipmentsForTOO() {
    return this.buildDefault('MoveWithNTSShipmentsForTOO');
  }

  /**
   * Use testharness to build HHG move with NTS for TOO
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithExternalNTSShipmentsForTOO() {
    return this.buildDefault('HHGMoveWithExternalNTSShipmentsForTOO');
  }

  /**
   * Use testharness to build HHG move with approved NTS shipment for TOO
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithApprovedNTSShipmentsForTOO() {
    return this.buildDefault('HHGMoveWithApprovedNTSShipmentsForTOO');
  }

  /**
   * Use testharness to build HHG move with NTS-R for TOO
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithNTSRShipmentsForTOO() {
    return this.buildDefault('HHGMoveWithNTSRShipmentsForTOO');
  }

  /**
   * Use testharness to build HHG move with approved NTS-R shipment for TOO
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithApprovedNTSRShipmentsForTOO() {
    return this.buildDefault('HHGMoveWithApprovedNTSRShipmentsForTOO');
  }

  /**
   * Use testharness to build HHG move with NTS-R for TOO
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithExternalNTSRShipmentsForTOO() {
    return this.buildDefault('HHGMoveWithExternalNTSRShipmentsForTOO');
  }

  /**
   * Use testharness to build hhg move for TIO
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithServiceItemsandPaymentRequestsForTIO() {
    return this.buildDefault('HHGMoveWithServiceItemsandPaymentRequestsForTIO');
  }

  /**
   * Use testharness to build hhg move needing SC approval
   * @returns {Promise<Move>}
   */
  async buildHHGMoveNeedsSC() {
    return this.buildDefault('HHGMoveNeedsSC');
  }

  /**
   * Use testharness to build hhg move needing SC approval
   * @returns {Promise<Move>}
   */
  async buildHHGMoveForSeparationNeedsSC() {
    return this.buildDefault('HHGMoveForSeparationNeedsSC');
  }

  /**
   * Use testharness to build hhg move needing SC approval
   * @returns {Promise<Move>}
   */
  async buildHHGMoveForRetireeNeedsSC() {
    return this.buildDefault('HHGMoveForRetireeNeedsSC');
  }

  /**
   * Use testharness to build hhg move with nts
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithNTSAndNeedsSC() {
    return this.buildDefault('HHGMoveWithNTSAndNeedsSC');
  }

  /**
   * Use testharness to build move with minimal NTS-R
   * @returns {Promise<Move>}
   */
  async buildMoveWithMinimalNTSRNeedsSC() {
    return this.buildDefault('MoveWithMinimalNTSRNeedsSC');
  }

  /**
   * Use testharness to build submitted move with ppm shipment for SC
   * @returns {Promise<Move>}
   */
  async buildSubmittedMoveWithPPMShipmentForSC() {
    return this.buildDefault('SubmittedMoveWithPPMShipmentForSC');
  }

  /**
   * Use testharness to build unsubmitted ppm move
   * @returns {Promise<Object>}
   */
  async buildUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights() {
    return this.buildDefault('UnSubmittedMoveWithPPMShipmentThroughEstimatedWeights');
  }

  /**
   * Use testharness to build unsubmitted ppm move
   * @returns {Promise<Object>}
   */
  async buildApprovedMoveWithPPMWithAboutFormComplete() {
    return this.buildDefault('ApprovedMoveWithPPMWithAboutFormComplete');
  }

  /**
   * Use testharness to build unsubmitted move with multiple ppm shipments
   * @returns {Promise<Object>}
   */
  async buildUnsubmittedMoveWithMultipleFullPPMShipmentComplete() {
    return this.buildDefault('UnsubmittedMoveWithMultipleFullPPMShipmentComplete');
  }

  /**
   * Use testharness to build approved move with ppm progear
   * @returns {Promise<Object>}
   */
  async buildApprovedMoveWithPPMProgearWeightTicket() {
    return this.buildDefault('ApprovedMoveWithPPMProgearWeightTicket');
  }

  /**
   * Use testharness to build approved move with ppm progear
   * @returns {Promise<Object>}
   */
  async buildApprovedMoveWithPPMWeightTicketOffice() {
    return this.buildDefault('ApprovedMoveWithPPMWeightTicketOffice');
  }

  /**
   * Use testharness to build approved move with ppm moving expenses
   * @returns {Promise<Object>}
   */
  async buildApprovedMoveWithPPMMovingExpense() {
    return this.buildDefault('ApprovedMoveWithPPMMovingExpense');
  }

  /**
   * Use testharness to build draft move with ppm departure date
   * @returns {Promise<Object>}
   */
  async buildDraftMoveWithPPMWithDepartureDate() {
    return this.buildDefault('DraftMoveWithPPMWithDepartureDate');
  }

  /**
   * Use testharness to build NTS-R move with payment request
   * @returns {Promise<Move>}
   */
  async buildNTSRMoveWithPaymentRequest() {
    return this.buildDefault('NTSRMoveWithPaymentRequest');
  }

  /**
   * Use testharness to build NTS-R move with service items payment request
   * @returns {Promise<Move>}
   */
  async buildNTSRMoveWithServiceItemsAndPaymentRequest() {
    return this.buildDefault('NTSRMoveWithServiceItemsAndPaymentRequest');
  }

  /**
   * @returns {Promise<Object>}
   */
  async buildPrimeSimulatorMoveNeedsShipmentUpdate() {
    return this.buildDefault('PrimeSimulatorMoveNeedsShipmentUpdate');
  }

  /**
   * @returns {Promise<Object>}
   */
  async buildWebhookSubscription() {
    return this.buildDefault('WebhookSubscription');
  }
}

export default TestHarness;
