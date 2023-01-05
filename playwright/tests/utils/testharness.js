// @ts-check
/**
 * Test Harness
 *
 * @param {import('@playwright/test').APIRequestContext} request
 */
export function newTestHarness(request) {
  /**
   * call the test harness build
   *
   * @param {string} action
   * @returns {Promise<Object>} Promise object for created data
   */
  const buildDefault = async (action) => {
    const r = await request.post(`/testharness/build/${action}`);
    if (!r.ok()) {
      const body = await r.body();
      throw Error(`Error with testharness build for '${action}': ${body}`);
    }
    const obj = /** @type {Object} */ await r.json();
    await r.dispose();
    return obj;
  };

  return {
    /**
     * @returns {Promise<Object>}
     */
    async buildDefaultAdminUser() {
      return buildDefault('DefaultAdminUser');
    },

    /**
     * build office user with TOO and TIO roles
     * @returns {Promise<Object>}
     */
    async buildOfficeUserWithTOOAndTIO() {
      return buildDefault('OfficeUserWithTOOAndTIO');
    },

    /**
     * @returns {Promise<Object>}
     */
    async buildDefaultMove() {
      return buildDefault('DefaultMove');
    },

    /**
     * @returns {Promise<Object>}
     */
    async buildMoveWithOrders() {
      return buildDefault('MoveWithOrders');
    },

    /**
     * @returns {Promise<Object>}
     */
    async buildSpouseProGearMove() {
      return buildDefault('SpouseProGearMove');
    },

    /**
     * @returns {Promise<Object>}
     */
    async buildNeedsOrdersUser() {
      return buildDefault('NeedsOrdersUser');
    },

    /**
     * @returns {Promise<Object>}
     */
    async buildPPMInProgressMove() {
      return buildDefault('PPMInProgressMove');
    },

    /**
     * @returns {Promise<Object>}
     */
    async buildWithShipmentMove() {
      return buildDefault('WithShipmentMove');
    },

    /**
     * Use testharness to build hhg move for TOO
     * @returns {Promise<Object>}
     */
    async buildHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO() {
      return buildDefault('HHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO');
    },

    /**
     * Use testharness to build retiree hhg move for TOO
     * @returns {Promise<Object>}
     */
    async buildHHGMoveWithRetireeForTOO() {
      return buildDefault('HHGMoveWithRetireeForTOO');
    },

    /**
     * Use testharness to build hhg move for TIO
     * @returns {Promise<Object>}
     */
    async buildHHGMoveWithServiceItemsandPaymentRequestsForTIO() {
      return buildDefault('HHGMoveWithServiceItemsandPaymentRequestsForTIO');
    },

    /**
     * Use testharness to build hhg move with nts
     * @returns {Promise<Object>}
     */
    async buildHHGMoveWithNTSAndNeedsSC() {
      return buildDefault('HHGMoveWithNTSAndNeedsSC');
    },

    /**
     * Use testharness to build NTS-R move with payment request
     * @returns {Promise<Object>}
     */
    async buildNTSRMoveWithPaymentRequest() {
      return buildDefault('NTSRMoveWithPaymentRequest');
    },

    /**
     * Use testharness to build NTS-R move with service items payment request
     * @returns {Promise<Object>}
     */
    async buildNTSRMoveWithServiceItemsAndPaymentRequest() {
      return buildDefault('NTSRMoveWithServiceItemsAndPaymentRequest');
    },

    /**
     * @returns {Promise<Object>}
     */
    async buildPrimeSimulatorMoveNeedsShipmentUpdate() {
      return buildDefault('PrimeSimulatorMoveNeedsShipmentUpdate');
    },

    /**
     * @returns {Promise<Object>}
     */
    async buildWebhookSubscription() {
      return buildDefault('WebhookSubscription');
    },
  };
}

export default newTestHarness;
