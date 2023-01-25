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
   * @returns {Promise<Object>}
   */
  async buildDefaultAdminUser() {
    return this.buildDefault('DefaultAdminUser');
  }

  /**
   * build office user with TOO and TIO roles
   * @returns {Promise<Object>}
   */
  async buildOfficeUserWithTOOAndTIO() {
    return this.buildDefault('OfficeUserWithTOOAndTIO');
  }

  /**
   * @returns {Promise<Object>}
   */
  async buildDefaultMove() {
    return this.buildDefault('DefaultMove');
  }

  /**
   * @returns {Promise<Object>}
   */
  async buildMoveWithOrders() {
    return this.buildDefault('MoveWithOrders');
  }

  /**
   * @returns {Promise<Object>}
   */
  async buildSpouseProGearMove() {
    return this.buildDefault('SpouseProGearMove');
  }

  /**
   * @returns {Promise<Object>}
   */
  async buildNeedsOrdersUser() {
    return this.buildDefault('NeedsOrdersUser');
  }

  /**
   * @returns {Promise<Object>}
   */
  async buildPPMInProgressMove() {
    return this.buildDefault('PPMInProgressMove');
  }

  /**
   * @returns {Promise<Object>}
   */
  async buildWithShipmentMove() {
    return this.buildDefault('WithShipmentMove');
  }

  /**
   * Use testharness to build hhg move for TOO
   * @returns {Promise<Object>}
   */
  async buildHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO() {
    return this.buildDefault('HHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO');
  }

  /**
   * Use testharness to build retiree hhg move for TOO
   * @returns {Promise<Object>}
   */
  async buildHHGMoveWithRetireeForTOO() {
    return this.buildDefault('HHGMoveWithRetireeForTOO');
  }

  /**
   * Use testharness to build HHG move with NTS for TOO
   * @returns {Promise<Object>}
   */
  async buildHHGMoveWithNTSShipmentsForTOO() {
    return this.buildDefault('HHGMoveWithNTSShipmentsForTOO');
  }

  /**
   * Use testharness to build move with NTS for TOO
   * @returns {Promise<Object>}
   */
  async buildMoveWithNTSShipmentsForTOO() {
    return this.buildDefault('MoveWithNTSShipmentsForTOO');
  }

  /**
   * Use testharness to build HHG move with NTS for TOO
   * @returns {Promise<Object>}
   */
  async buildHHGMoveWithExternalNTSShipmentsForTOO() {
    return this.buildDefault('HHGMoveWithExternalNTSShipmentsForTOO');
  }

  /**
   * Use testharness to build HHG move with approved NTS shipment for TOO
   * @returns {Promise<Object>}
   */
  async buildHHGMoveWithApprovedNTSShipmentsForTOO() {
    return this.buildDefault('HHGMoveWithApprovedNTSShipmentsForTOO');
  }

  /**
   * Use testharness to build HHG move with NTS-R for TOO
   * @returns {Promise<Object>}
   */
  async buildHHGMoveWithNTSRShipmentsForTOO() {
    return this.buildDefault('HHGMoveWithNTSRShipmentsForTOO');
  }

  /**
   * Use testharness to build HHG move with approved NTS-R shipment for TOO
   * @returns {Promise<Object>}
   */
  async buildHHGMoveWithApprovedNTSRShipmentsForTOO() {
    return this.buildDefault('HHGMoveWithApprovedNTSRShipmentsForTOO');
  }

  /**
   * Use testharness to build HHG move with NTS-R for TOO
   * @returns {Promise<Object>}
   */
  async buildHHGMoveWithExternalNTSRShipmentsForTOO() {
    return this.buildDefault('HHGMoveWithExternalNTSRShipmentsForTOO');
  }

  /**
   * Use testharness to build hhg move for TIO
   * @returns {Promise<Object>}
   */
  async buildHHGMoveWithServiceItemsandPaymentRequestsForTIO() {
    return this.buildDefault('HHGMoveWithServiceItemsandPaymentRequestsForTIO');
  }

  /**
   * Use testharness to build hhg move with nts
   * @returns {Promise<Object>}
   */
  async buildHHGMoveWithNTSAndNeedsSC() {
    return this.buildDefault('HHGMoveWithNTSAndNeedsSC');
  }

  /**
   * Use testharness to build NTS-R move with payment request
   * @returns {Promise<Object>}
   */
  async buildNTSRMoveWithPaymentRequest() {
    return this.buildDefault('NTSRMoveWithPaymentRequest');
  }

  /**
   * Use testharness to build NTS-R move with service items payment request
   * @returns {Promise<Object>}
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
