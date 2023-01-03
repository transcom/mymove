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
     * Use testharness to build hhg move
     * @returns {Promise<Object>}
     */
    async buildHHGMoveWithServiceItemsAndPaymentRequestsAndFiles() {
      return buildDefault('HHGMoveWithServiceItemsAndPaymentRequestsAndFiles');
    },

    /**
     * Use testharness to build hhg move with nts
     * @returns {Promise<Object>}
     */
    async buildHHGMoveWithNTSAndNeedsSC() {
      return buildDefault('HHGMoveWithNTSAndNeedsSC');
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
