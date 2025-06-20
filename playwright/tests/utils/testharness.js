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
   * @property {string} okta_email
   */

  /**
   * @typedef {Object} Admin
   * @property {string} id
   * @property {string} user_id
   */

  /**
   * @typedef {Object} Move
   * @property {string} id
   * @property {string} locator
   * @property {Object} Orders
   * @property {Object} OriginDutyLocation
   * @property {Object} OriginDutyLocation.Address
   * @property {string} OriginDutyLocation.Address.PostalCode
   * @property {Object} Orders.NewDutyLocation
   * @property {string} Orders.NewDutyLocation.name
   * @property {Object} Orders.ServiceMember
   * @property {string} Orders.ServiceMember.edipi
   * @property {string} Orders.ServiceMember.last_name
   * @property {string} Orders.ServiceMember.user_id
   * @property {Object} CloseoutOffice
   * @property {string} CloseoutOffice.name
   * @property {Array.<Object>} MTOShipments
   * @property {string} MTOShipments.ID
   * @property {Array.<Object>} MTOServiceItems
   * @property {string} MTOServiceItems.ID
   * @property {Object} MTOServiceItems.ReService
   * @property {string} MTOServiceItems.ReService.ID
   * @property {string} MTOServiceItems.ReService.Code
   *
   */

  /**
   * @typedef {Object} TransportationAccountingCode
   * @property {string} id
   * @property {string} TAC
   * @property {string} LoaSysID
   */

  /**
   * @typedef {Object} WebhookSubscription
   * @property {string} ID
   * @property {string} SubscriberID
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
   * @returns {Promise<Admin>}
   */
  async buildDefaultSuperAdminUser() {
    return this.buildDefault('SuperAdminUser');
  }

  /**
   * @returns {Promise<User>}
   */
  async buildRequestedOfficeUser() {
    return this.buildDefault('RequestedOfficeUser');
  }

  /**
   * build requested office user with a privilege request
   * @returns {Promise<User>}
   */
  async buildRequestedOfficeUserWithPrivilege() {
    return this.buildDefault('RequestedOfficeUserWithPrivilege');
  }

  /**
   * @returns {Promise<User>}
   */
  async buildRejectedOfficeUser() {
    return this.buildDefault('RejectedOfficeUser');
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
  async buildPartialPPMMoveReadyForCloseout() {
    return this.buildDefault('HHGMoveWithPPMShipmentsReadyForCloseout');
  }

  /**
   * @returns {Promise<Move>}
   */
  async buildPartialPPMMoveReadyForCounseling() {
    return this.buildDefault('HHGMoveWithPPMShipmentsReadyForCounseling');
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
   * Use testharness to build a move with an hhg shipment in SIT
   * @returns {Promise<Move>}
   */
  async buildHHGMoveInSIT() {
    return this.buildDefault('HHGMoveInSIT');
  }

  /**
   * Use testharness to build a move with an hhg shipment with a past origin and destination SIT
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithPastSITs() {
    return this.buildDefault('HHGMoveWithPastSITs');
  }

  /**
   *    * Use testharness to build a move with an hhg shipment in SIT without delivery address
   * @returns {Promise<Move>}
   */

  async buildHHGMoveInSITNoDestinationSITOutDate() {
    return this.buildDefault('HHGMoveInSITNoDestinationSITOutDate');
  }

  /**
   * Use testharness to build a move with an hhg shipment in SIT without excess weight
   * @returns {Promise<Move>}
   */
  async buildHHGMoveInSITNoExcessWeight() {
    return this.buildDefault('HHGMoveInSITNoExcessWeight');
  }

  /**
   * Use testharness to build a move with an hhg shipment in SIT and a pending SIT extension
   * @returns {Promise<Move>}
   */
  async buildHHGMoveInSITWithPendingExtension() {
    return this.buildDefault('HHGMoveInSITWithPendingExtension');
  }

  /**
   * Use testharness to build a move with an hhg shipment in SIT with an allowance that ends today
   * @returns {Promise<Move>}
   */
  async buildHHGMoveInSITEndsToday() {
    return this.buildDefault('HHGMoveInSITEndsToday');
  }

  /**
   * Use testharness to build a move with an hhg shipment in SIT with an allowance that ends tomorrow
   * @returns {Promise<Move>}
   */
  async buildHHGMoveInSITEndsTomorrow() {
    return this.buildDefault('HHGMoveInSITEndsTomorrow');
  }

  /**
   * Use testharness to build a move with an hhg shipment in SIT with an allowance that ended yesterday
   * @returns {Promise<Move>}
   */
  async buildHHGMoveInSITEndsYesterday() {
    return this.buildDefault('HHGMoveInSITEndsYesterday');
  }

  /**
   * Use testharness to build a move with an hhg shipment in SIT that departed storage before the allowance
   * was exhausted
   * @returns {Promise<Move>}
   */
  async buildHHGMoveInSITDeparted() {
    return this.buildDefault('HHGMoveInSITDeparted');
  }

  /**
   * Use testharness to build a move with an hhg shipment that hasn't yet entered SIT
   * @returns {Promise<Move>}
   */
  async buildHHGMoveInSITStartsInFuture() {
    return this.buildDefault('HHGMoveInSITStartsInFuture');
  }

  /**
   * Use testharness to build a move with an hhg shipment with SIT that hasn't yet been approved
   * @returns {Promise<Move>}
   */
  async buildHHGMoveInSITNotApproved() {
    return this.buildDefault('HHGMoveInSITNotApproved');
  }

  /**
   * Use testharness to build hhg move for TOO
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO() {
    return this.buildDefault('HHGMoveWithServiceItemsAndPaymentRequestsAndFilesForTOO');
  }

  /**
   * Use testharness to build terminated hhg move for TOO
   * They should still be able to approve/deny service items
   * but not request cancellation, diversion, reweigh, or other shipment-specific
   * actions. Only actions associated with the shipment
   * @returns {Promise<Move>}
   */
  async buildHHGMoveInTerminatedStatus() {
    return this.buildDefault('HHGMoveInTerminatedStatus');
  }

  /**
   * Use testharness to build hhg move for TOO with Alaska address
   * @returns {Promise<Move>}
   */
  async buildInternationalAlaskaBasicHHGMoveForTOO() {
    return this.buildDefault('InternationalAlaskaBasicHHGMoveForTOO');
  }

  /**
   * Use testharness to build hhg move with international crating service items for TOO
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithIntlCratingServiceItemsTOO() {
    return this.buildDefault('HHGMoveWithIntlCratingServiceItemsTOO');
  }

  /**
   * Use testharness to build hhg move with international crating service items for TOO
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithIntlShuttleServiceItemsTOO() {
    return this.buildDefault('HHGMoveWithIntlShuttleServiceItemsTOO');
  }

  /**
   * Use testharness to build hhg move for TOO with actualPickupDate in the past
   * @returns {Promise<Move>}
   */
  async buildHHGMoveForTOOAfterActualPickupDate() {
    return this.buildDefault('HHGMoveForTOOAfterActualPickupDate');
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
   * Use testharness to build iHHG move with iNTS basic accessory service items
   * for TIO review
   * eg, INPK
   * @returns {Promise<Move>}
   */
  async buildInternationalHHGIntoInternationalNTSMoveWithServiceItemsandPaymentRequestsForTIO() {
    return this.buildDefault('InternationalHHGIntoInternationalNTSMoveWithServiceItemsandPaymentRequestsForTIO');
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
   * Use testharness to build hhg move for TIO
   * @returns {Promise<Move>}
   */
  async buildnternationalHHGMoveWithServiceItemsandPaymentRequestsForTIO() {
    return this.buildDefault('InternationalHHGMoveWithServiceItemsandPaymentRequestsForTIO');
  }

  /**
   * Use testharness to build ihhg move for TIO
   * @returns {Promise<Move>}
   */
  async buildIntlHHGMoveWithCratingUncratingServiceItemsAndPaymentRequestsForTIO() {
    return this.buildDefault('IntlHHGMoveWithCratingUncratingServiceItemsAndPaymentRequestsForTIO');
  }

  /**
   * Use testharness to build hhg move for QAE
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithServiceItemsandPaymentRequestReviewedForQAE() {
    return this.buildDefault('HHGMoveWithServiceItemsandPaymentRequestReviewedForQAE');
  }

  /**
   * Use testharness to build hhg move for QAE
   * @returns {Promise<Move>}
   */
  async buildHHGMoveWithServiceItemsandPaymentRequestWithDocsReviewedForQAE() {
    return this.buildDefault('HHGMoveWithServiceItemsandPaymentRequestWithDocsReviewedForQAE');
  }

  /**
   * Use testharness to build hhg move needing SC approval
   * @returns {Promise<Move>}
   */
  async buildHHGMoveNeedsSC() {
    return this.buildDefault('HHGMoveNeedsSC');
  }

  /**
   * Use testharness to build international hhg move needing SC approval
   * @returns {Promise<Move>}
   */
  async buildIntlHHGMoveNeedsSC() {
    return this.buildDefault('IntlHHGMoveNeedsSC');
  }

  /**
   * Use testharness to build hhg move needing SC approval in a non-default GBLOC
   * @returns {Promise<Move>}
   */
  async buildHHGMoveNeedsSCInOtherGBLOC() {
    return this.buildDefault('HHGMoveNeedsSCOtherGBLOC');
  }

  /**
   * Use testharness to build hhg move as USMC needing SC approval
   * @returns {Promise<Move>}
   */
  async buildHHGMoveAsUSMCNeedsSC() {
    return this.buildDefault('HHGMoveAsUSMCNeedsSC');
  }

  /**
   * Use testharness to build hhg move with amended orders
   * @returns {Promise<Move>}
   */
  async buildHHGWithAmendedOrders() {
    return this.buildDefault('HHGMoveWithAmendedOrders');
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
   * Use testharness to build Mobile move
   * @returns {Promise<Move>}
   */
  async buildMobileHomeMoveNeedsSC() {
    return this.buildDefault('MobileHomeMoveNeedsSC');
  }

  /**
   * Use testharness to build a good TAC and LOA combination, return the TAC
   * so that office users can input the TAC, and preview the LOA (If the
   * form params are good for the lookup. Eg, service member branch,
   * effective date, that sort)
   * @returns {Promise<TransportationAccountingCode>}
   */
  async buildGoodTACAndLoaCombination() {
    return this.buildDefault('GoodTACAndLoaCombination');
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
   * Use testharness to build submitted move with too created ppm shipment for SC
   * @returns {Promise<Move>}
   */
  async buildApprovedMoveWithSubmittedPPMShipmentForSC() {
    return this.buildDefault('ApprovedMoveWithSubmittedPPMShipmentForSC');
  }

  /**
   * Use testharness to build unsubmitted ppm move
   * @returns {Promise<Move>}
   */
  async buildUnSubmittedMoveWithPPMShipmentThroughEstimatedWeights() {
    return this.buildDefault('UnSubmittedMoveWithPPMShipmentThroughEstimatedWeights');
  }

  /**
   * Use testharness to build unsubmitted ppm move
   * @returns {Promise<Move>}
   */
  async buildApprovedMoveWithPPMWithAboutFormComplete() {
    return this.buildDefault('ApprovedMoveWithPPMWithAboutFormComplete');
  }

  /**
   * Use testharness to build unsubmitted move with multiple ppm shipments
   * @returns {Promise<Move>}
   */
  async buildUnsubmittedMoveWithMultipleFullPPMShipmentComplete() {
    return this.buildDefault('UnsubmittedMoveWithMultipleFullPPMShipmentComplete');
  }

  /**
   * Use testharness to build approved move with ppm progear
   * @returns {Promise<Move>}
   */
  async buildApprovedMoveWithPPMProgearWeightTicket() {
    return this.buildDefault('ApprovedMoveWithPPMProgearWeightTicket');
  }

  /**
   * Use testharness to build Use testharness to build submitted move with ppm and pro-gear
   * @returns {Promise<Move>}
   */
  async buildApprovedMoveWithPPMProgearWeightTicketOffice() {
    return this.buildDefault('ApprovedMoveWithPPMProgearWeightTicketOffice');
  }

  /**
   * Use testharness to build Use testharness to build submitted move with ppm and pro-gear
   * @returns {Promise<Move>}
   */
  async buildApprovedMoveWithPPMWithMultipleProgearWeightTicketsOffice() {
    return this.buildDefault('ApprovedMoveWithPPMWithMultipleProgearWeightTicketsOffice');
  }

  /**
   * Use testharness to build Use testharness to build submitted move with ppm and pro-gear
   * @returns {Promise<Move>}
   */
  async buildApprovedMoveWithPPMWithMultipleProgearWeightTicketsOffice2() {
    return this.buildDefault('ApprovedMoveWithPPMWithMultipleProgearWeightTicketsOffice2');
  }

  /**
   * Use testharness to build Use testharness to build submitted move with ppm and pro-gear - civilian
   * @returns {Promise<Move>}
   */
  async buildApprovedMoveWithPPMProgearWeightTicketOfficeCivilian() {
    return this.buildDefault('ApprovedMoveWithPPMProgearWeightTicketOfficeCivilian');
  }

  /**
   * Use testharness to build submitted move with ppm and weight ticket
   * @returns {Promise<Object>}
   */
  async buildApprovedMoveWithPPMWeightTicketOffice() {
    return this.buildDefault('ApprovedMoveWithPPMWeightTicketOffice');
  }

  /**
   * Use testharness to build submitted move with partial ppm and weight ticket
   * @returns {Promise<Object>}
   */
  async buildApprovedMoveWithPPMWeightTicketOfficeWithHHG() {
    return this.buildDefault('ApprovedMoveWithPPMWeightTicketOfficeWithHHG');
  }

  /**
   * Use testharness to build approved move with ppm moving expenses
   * @returns {Promise<Move>}
   */
  async buildApprovedMoveWithPPMMovingExpense() {
    return this.buildDefault('ApprovedMoveWithPPMMovingExpense');
  }

  /**
   * Use testharness to build submitted move with ppm and moving expense
   * @returns {Promise<Object>}
   */
  async buildApprovedMoveWithPPMMovingExpenseOffice() {
    return this.buildDefault('ApprovedMoveWithPPMMovingExpenseOffice');
  }

  /**
   * Use testharness to build submitted move with ppm and all doc types
   * @returns {Promise<Object>}
   */
  async buildApprovedMoveWithPPMAllDocTypesOffice() {
    return this.buildDefault('ApprovedMoveWithPPMAllDocTypesOffice');
  }

  /**
   * Use testharness to build approved move with two ppm shipments and excess weight
   * @returns {Promise<Move>}
   */
  async buildApprovedMoveWithPPMShipmentAndExcessWeight() {
    return this.buildDefault('ApprovedMoveWithPPMShipmentAndExcessWeight');
  }

  /**
   * Use testharness to build draft move with ppm departure date
   * @returns {Promise<Move>}
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
   * @returns {Promise<Move>}
   */
  async buildPrimeSimulatorMoveNeedsShipmentUpdate() {
    return this.buildDefault('PrimeSimulatorMoveNeedsShipmentUpdate');
  }

  /**
   * @returns {Promise<WebhookSubscription>}
   */
  async buildWebhookSubscription() {
    return this.buildDefault('WebhookSubscription');
  }
  /**
   *
   * @returns {Promise<Object>}
   */

  async buildHHGMoveWithAddressChangeRequest() {
    return this.buildDefault('HHGMoveWithAddressChangeRequest');
  }

  /**
   *
   * @returns {Promise<Object>}
   */

  async buildNTSRMoveWithAddressChangeRequest() {
    return this.buildDefault('NTSRMoveWithAddressChangeRequest');
  }

  /**
   * Use testharness to build boat move needing SC
   * @returns {Promise<Move>}
   */

  async buildBoatHaulAwayMoveNeedsSC() {
    return this.buildDefault('BoatHaulAwayMoveNeedsSC');
  }

  /**
   * Use testharness to build boat move needing TOO approval
   * @returns {Promise<Move>}
   */

  async buildBoatHaulAwayMoveNeedsTOOApproval() {
    return this.buildDefault('BoatHaulAwayMoveNeedsTOOApproval');
  }

  /**
   * build office user with Customer role
   * @returns {Promise<User>}
   */
  async buildOfficeUserWithCustomer() {
    return this.buildDefault('OfficeUserWithCustomer');
  }

  /**
   * build office user with Contracting Officer role
   * @returns {Promise<User>}
   */
  async buildOfficeUserWithContractingOfficer() {
    return this.buildDefault('OfficeUserWithContractingOfficer');
  }

  /**
   * build office user with Prime Simulator role
   * @returns {Promise<User>}
   */
  async buildOfficeUserWithPrimeSimulator() {
    return this.buildDefault('OfficeUserWithPrimeSimulator');
  }

  /**
   * build office user with GSR role
   * @returns {Promise<User>}
   */
  async buildOfficeUserWithGSR() {
    return this.buildDefault('OfficeUserWithGSR');
  }

  /**
   * Use testharness to build international move with requested origin SIT
   * @returns {Promise<Move>}
   */
  async buildIntlHHGMoveOriginSITRequestedAKZone1Army() {
    return this.buildDefault('IntlHHGMoveOriginSITRequestedAKZone1Army');
  }

  async buildIntlHHGMoveOriginSITRequestedAKZone2Army() {
    return this.buildDefault('IntlHHGMoveOriginSITRequestedAKZone2Army');
  }

  async buildIntlHHGMoveOriginSITRequestedAKZone1AirForce() {
    return this.buildDefault('IntlHHGMoveOriginSITRequestedAKZone1AirForce');
  }

  async buildIntlHHGMoveOriginSITRequestedAKZone2AirForce() {
    return this.buildDefault('IntlHHGMoveOriginSITRequestedAKZone2AirForce');
  }

  async buildIntlHHGMoveOriginSITRequestedAKZone1SpaceForce() {
    return this.buildDefault('IntlHHGMoveOriginSITRequestedAKZone1SpaceForce');
  }

  async buildIntlHHGMoveOriginSITRequestedAKZone2SpaceForce() {
    return this.buildDefault('IntlHHGMoveOriginSITRequestedAKZone2SpaceForce');
  }

  async buildIntlHHGMoveOriginSITRequestedAKZone1USMC() {
    return this.buildDefault('IntlHHGMoveOriginSITRequestedAKZone1USMC');
  }

  async buildIntlHHGMoveOriginSITRequestedAKZone2USMC() {
    return this.buildDefault('IntlHHGMoveOriginSITRequestedAKZone2USMC');
  }

  /**
   * Use testharness to build international move with requested destination SIT
   * @returns {Promise<Move>}
   */
  async buildIntlHHGMoveDestSITRequestedAKZone1Army() {
    return this.buildDefault('IntlHHGMoveDestSITRequestedAKZone1Army');
  }

  async buildIntlHHGMoveDestSITRequestedAKZone2Army() {
    return this.buildDefault('IntlHHGMoveDestSITRequestedAKZone2Army');
  }

  async buildIntlHHGMoveDestSITRequestedAKZone1AirForce() {
    return this.buildDefault('IntlHHGMoveDestSITRequestedAKZone1AirForce');
  }

  async buildIntlHHGMoveDestSITRequestedAKZone2AirForce() {
    return this.buildDefault('IntlHHGMoveDestSITRequestedAKZone2AirForce');
  }

  async buildIntlHHGMoveDestSITRequestedAKZone1SpaceForce() {
    return this.buildDefault('IntlHHGMoveDestSITRequestedAKZone1SpaceForce');
  }

  async buildIntlHHGMoveDestSITRequestedAKZone2SpaceForce() {
    return this.buildDefault('IntlHHGMoveDestSITRequestedAKZone2SpaceForce');
  }

  async buildIntlHHGMoveDestSITRequestedAKZone1USMC() {
    return this.buildDefault('IntlHHGMoveDestSITRequestedAKZone1USMC');
  }

  async buildIntlHHGMoveDestSITRequestedAKZone2USMC() {
    return this.buildDefault('IntlHHGMoveDestSITRequestedAKZone2USMC');
  }

  /**
   * Use testharness to build international move with both requested origin & destination SIT
   * @returns {Promise<Move>}
   */
  async buildIntlHHGMoveBothSITRequestedAKZone1Army() {
    return this.buildDefault('IntlHHGMoveBothSITRequestedAKZone1Army');
  }

  async buildIntlHHGMoveBothSITRequestedAKZone2Army() {
    return this.buildDefault('IntlHHGMoveBothSITRequestedAKZone2Army');
  }

  async buildIntlHHGMoveBothSITRequestedAKZone1AirForce() {
    return this.buildDefault('IntlHHGMoveBothSITRequestedAKZone1AirForce');
  }

  async buildIntlHHGMoveBothSITRequestedAKZone2AirForce() {
    return this.buildDefault('IntlHHGMoveBothSITRequestedAKZone2AirForce');
  }

  async buildIntlHHGMoveBothSITRequestedAKZone1SpaceForce() {
    return this.buildDefault('IntlHHGMoveBothSITRequestedAKZone1SpaceForce');
  }

  async buildIntlHHGMoveBothSITRequestedAKZone2SpaceForce() {
    return this.buildDefault('IntlHHGMoveBothSITRequestedAKZone2SpaceForce');
  }

  async buildIntlHHGMoveBothSITRequestedAKZone1USMC() {
    return this.buildDefault('IntlHHGMoveBothSITRequestedAKZone1USMC');
  }

  async buildIntlHHGMoveBothSITRequestedAKZone2USMC() {
    return this.buildDefault('IntlHHGMoveBothSITRequestedAKZone2USMC');
  }

  /**
   * Use testharness to build international move with requested origin shuttle service
   * @returns {Promise<Move>}
   */
  async buildIntlHHGMoveOriginShuttleRequestedAKZone1Army() {
    return this.buildDefault('IntlHHGMoveOriginShuttleRequestedAKZone1Army');
  }

  async buildIntlHHGMoveOriginShuttleRequestedAKZone2Army() {
    return this.buildDefault('IntlHHGMoveOriginShuttleRequestedAKZone2Army');
  }

  async buildIntlHHGMoveOriginShuttleRequestedAKZone1AirForce() {
    return this.buildDefault('IntlHHGMoveOriginShuttleRequestedAKZone1AirForce');
  }

  async buildIntlHHGMoveOriginShuttleRequestedAKZone2AirForce() {
    return this.buildDefault('IntlHHGMoveOriginShuttleRequestedAKZone2AirForce');
  }

  async buildIntlHHGMoveOriginShuttleRequestedAKZone1SpaceForce() {
    return this.buildDefault('IntlHHGMoveOriginShuttleRequestedAKZone1SpaceForce');
  }

  async buildIntlHHGMoveOriginShuttleRequestedAKZone2SpaceForce() {
    return this.buildDefault('IntlHHGMoveOriginShuttleRequestedAKZone2SpaceForce');
  }

  async buildIntlHHGMoveOriginShuttleRequestedAKZone1USMC() {
    return this.buildDefault('IntlHHGMoveOriginShuttleRequestedAKZone1USMC');
  }

  async buildIntlHHGMoveOriginShuttleRequestedAKZone2USMC() {
    return this.buildDefault('IntlHHGMoveOriginShuttleRequestedAKZone2USMC');
  }

  /**
   * Use testharness to build international move with requested destination shuttle service
   * @returns {Promise<Move>}
   */
  async buildIntlHHGMoveDestShuttleRequestedAKZone1Army() {
    return this.buildDefault('IntlHHGMoveDestShuttleRequestedAKZone1Army');
  }

  async buildIntlHHGMoveDestShuttleRequestedAKZone2Army() {
    return this.buildDefault('IntlHHGMoveDestShuttleRequestedAKZone2Army');
  }

  async buildIntlHHGMoveDestShuttleRequestedAKZone1AirForce() {
    return this.buildDefault('IntlHHGMoveDestShuttleRequestedAKZone1AirForce');
  }

  async buildIntlHHGMoveDestShuttleRequestedAKZone2AirForce() {
    return this.buildDefault('IntlHHGMoveDestShuttleRequestedAKZone2AirForce');
  }

  async buildIntlHHGMoveDestShuttleRequestedAKZone1SpaceForce() {
    return this.buildDefault('IntlHHGMoveDestShuttleRequestedAKZone1SpaceForce');
  }

  async buildIntlHHGMoveDestShuttleRequestedAKZone2SpaceForce() {
    return this.buildDefault('IntlHHGMoveDestShuttleRequestedAKZone2SpaceForce');
  }

  async buildIntlHHGMoveDestShuttleRequestedAKZone1USMC() {
    return this.buildDefault('IntlHHGMoveDestShuttleRequestedAKZone1USMC');
  }

  async buildIntlHHGMoveDestShuttleRequestedAKZone2USMC() {
    return this.buildDefault('IntlHHGMoveDestShuttleRequestedAKZone2USMC');
  }

  /**
   * Use testharness to build international move with both requested origin & destination shuttle service
   * @returns {Promise<Move>}
   */
  async buildIntlHHGMoveBothShuttleRequestedAKZone1Army() {
    return this.buildDefault('IntlHHGMoveBothShuttleRequestedAKZone1Army');
  }

  async buildIntlHHGMoveBothShuttleRequestedAKZone2Army() {
    return this.buildDefault('IntlHHGMoveBothShuttleRequestedAKZone2Army');
  }

  async buildIntlHHGMoveBothShuttleRequestedAKZone1AirForce() {
    return this.buildDefault('IntlHHGMoveBothShuttleRequestedAKZone1AirForce');
  }

  async buildIntlHHGMoveBothShuttleRequestedAKZone2AirForce() {
    return this.buildDefault('IntlHHGMoveBothShuttleRequestedAKZone2AirForce');
  }

  async buildIntlHHGMoveBothShuttleRequestedAKZone1SpaceForce() {
    return this.buildDefault('IntlHHGMoveBothShuttleRequestedAKZone1SpaceForce');
  }

  async buildIntlHHGMoveBothShuttleRequestedAKZone2SpaceForce() {
    return this.buildDefault('IntlHHGMoveBothShuttleRequestedAKZone2SpaceForce');
  }

  async buildIntlHHGMoveBothShuttleRequestedAKZone1USMC() {
    return this.buildDefault('IntlHHGMoveBothShuttleRequestedAKZone1USMC');
  }

  async buildIntlHHGMoveBothShuttleRequestedAKZone2USMC() {
    return this.buildDefault('IntlHHGMoveBothShuttleRequestedAKZone2USMC');
  }

  /**
   * Use testharness to build international move with requested destination address request
   * @returns {Promise<Move>}
   */
  async buildIntlHHGMoveDestAddressRequestedAKZone1Army() {
    return this.buildDefault('IntlHHGMoveDestAddressRequestedAKZone1Army');
  }

  async buildIntlHHGMoveDestAddressRequestedAKZone2Army() {
    return this.buildDefault('IntlHHGMoveDestAddressRequestedAKZone2Army');
  }

  async buildIntlHHGMoveDestAddressRequestedAKZone1AirForce() {
    return this.buildDefault('IntlHHGMoveDestAddressRequestedAKZone1AirForce');
  }

  async buildIntlHHGMoveDestAddressRequestedAKZone2AirForce() {
    return this.buildDefault('IntlHHGMoveDestAddressRequestedAKZone2AirForce');
  }

  async buildIntlHHGMoveDestAddressRequestedAKZone1SpaceForce() {
    return this.buildDefault('IntlHHGMoveDestAddressRequestedAKZone1SpaceForce');
  }

  async buildIntlHHGMoveDestAddressRequestedAKZone2SpaceForce() {
    return this.buildDefault('IntlHHGMoveDestAddressRequestedAKZone2SpaceForce');
  }

  async buildIntlHHGMoveDestAddressRequestedAKZone1USMC() {
    return this.buildDefault('IntlHHGMoveDestAddressRequestedAKZone1USMC');
  }

  async buildIntlHHGMoveDestAddressRequestedAKZone2USMC() {
    return this.buildDefault('IntlHHGMoveDestAddressRequestedAKZone2USMC');
  }

  /**
   * Use testharness to build international move with a pending SIT extension request with origin SIT
   * @returns {Promise<Move>}
   */
  async buildIntlHHGMoveOriginSITExtensionRequestedAKZone1Army() {
    return this.buildDefault('IntlHHGMoveOriginSITExtensionRequestedAKZone1Army');
  }

  async buildIntlHHGMoveOriginSITExtensionRequestedAKZone2Army() {
    return this.buildDefault('IntlHHGMoveOriginSITExtensionRequestedAKZone2Army');
  }

  async buildIntlHHGMoveOriginSITExtensionRequestedAKZone1AirForce() {
    return this.buildDefault('IntlHHGMoveOriginSITExtensionRequestedAKZone1AirForce');
  }

  async buildIntlHHGMoveOriginSITExtensionRequestedAKZone2AirForce() {
    return this.buildDefault('IntlHHGMoveOriginSITExtensionRequestedAKZone2AirForce');
  }

  async buildIntlHHGMoveOriginSITExtensionRequestedAKZone1SpaceForce() {
    return this.buildDefault('IntlHHGMoveOriginSITExtensionRequestedAKZone1SpaceForce');
  }

  async buildIntlHHGMoveOriginSITExtensionRequestedAKZone2SpaceForce() {
    return this.buildDefault('IntlHHGMoveOriginSITExtensionRequestedAKZone2SpaceForce');
  }

  async buildIntlHHGMoveOriginSITExtensionRequestedAKZone1USMC() {
    return this.buildDefault('IntlHHGMoveOriginSITExtensionRequestedAKZone1USMC');
  }

  async buildIntlHHGMoveOriginSITExtensionRequestedAKZone2USMC() {
    return this.buildDefault('IntlHHGMoveOriginSITExtensionRequestedAKZone2USMC');
  }

  /**
   * Use testharness to build international move with a pending SIT extension request with destination SIT
   * @returns {Promise<Move>}
   */
  async buildIntlHHGMoveDestSITExtensionRequestedAKZone1Army() {
    return this.buildDefault('IntlHHGMoveDestSITExtensionRequestedAKZone1Army');
  }

  async buildIntlHHGMoveDestSITExtensionRequestedAKZone2Army() {
    return this.buildDefault('IntlHHGMoveDestSITExtensionRequestedAKZone2Army');
  }

  async buildIntlHHGMoveDestSITExtensionRequestedAKZone1AirForce() {
    return this.buildDefault('IntlHHGMoveDestSITExtensionRequestedAKZone1AirForce');
  }

  async buildIntlHHGMoveDestSITExtensionRequestedAKZone2AirForce() {
    return this.buildDefault('IntlHHGMoveDestSITExtensionRequestedAKZone2AirForce');
  }

  async buildIntlHHGMoveDestSITExtensionRequestedAKZone1SpaceForce() {
    return this.buildDefault('IntlHHGMoveDestSITExtensionRequestedAKZone1SpaceForce');
  }

  async buildIntlHHGMoveDestSITExtensionRequestedAKZone2SpaceForce() {
    return this.buildDefault('IntlHHGMoveDestSITExtensionRequestedAKZone2SpaceForce');
  }

  async buildIntlHHGMoveDestSITExtensionRequestedAKZone1USMC() {
    return this.buildDefault('IntlHHGMoveDestSITExtensionRequestedAKZone1USMC');
  }

  async buildIntlHHGMoveDestSITExtensionRequestedAKZone2USMC() {
    return this.buildDefault('IntlHHGMoveDestSITExtensionRequestedAKZone2USMC');
  }

  /**
   * Use testharness to build international move with an at risk of excess weight
   * @returns {Promise<Move>}
   */
  async buildIntlHHGMoveExcessWeightAKZone1Army() {
    return this.buildDefault('IntlHHGMoveExcessWeightAKZone1Army');
  }

  async buildIntlHHGMoveExcessWeightAKZone2Army() {
    return this.buildDefault('IntlHHGMoveExcessWeightAKZone2Army');
  }

  /**
   * Use testharness to build international move with an at risk of UB excess weight
   * @returns {Promise<Move>}
   */
  async buildIntlUBMoveExcessWeightAKZone1Army() {
    return this.buildDefault('IntlUBMoveExcessWeightAKZone1Army');
  }

  async buildIntlUBMoveExcessWeightAKZone2Army() {
    return this.buildDefault('IntlUBMoveExcessWeightAKZone2Army');
  }
}
export default TestHarness;
