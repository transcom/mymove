/* eslint-disable react/jsx-props-no-spreading */
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { func, arrayOf } from 'prop-types';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';
import { generatePath } from 'react-router-dom';

import { isBooleanFlagEnabled } from '../../utils/featureFlags';
import { FEATURE_FLAG_KEYS, SHIPMENT_OPTIONS } from '../../shared/constants';

import ConnectedMoveInfoModal from 'components/Customer/modals/MoveInfoModal/MoveInfoModal';
import ConnectedStorageInfoModal from 'components/Customer/modals/StorageInfoModal/StorageInfoModal';
import ConnectedBoatInfoModal from 'components/Customer/modals/BoatInfoModal/BoatInfoModal';
import SelectableCard from 'components/Customer/SelectableCard';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import NotificationScrollToTop from 'components/NotificationScrollToTop';
import { generalRoutes, customerRoutes } from 'constants/routes';
import styles from 'pages/MyMove/SelectShipmentType.module.scss';
import { loadMTOShipments as loadMTOShipmentsAction } from 'shared/Entities/modules/mtoShipments';
import { updateMove as updateMoveAction } from 'store/entities/actions';
import { selectMTOShipmentsForCurrentMove } from 'store/entities/selectors';
import formStyles from 'styles/form.module.scss';
import { MoveTaskOrderShape } from 'types/order';
import { ShipmentShape } from 'types/shipment';
import determineShipmentInfo from 'utils/shipmentInfo';
import withRouter from 'utils/routing';
import { RouterShape } from 'types';
import { selectMove } from 'shared/Entities/modules/moves';

export class SelectShipmentType extends Component {
  constructor(props) {
    super(props);
    this.state = {
      showStorageInfoModal: false,
      showMoveInfoModal: false,
      showBoatInfoModal: false,
      errorMessage: null,
      enablePPM: false,
      enableNTS: false,
      enableNTSR: false,
      enableBoat: false,
    };
  }

  componentDidMount() {
    const { loadMTOShipments, move } = this.props;
    loadMTOShipments(move.id);
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.PPM).then((enabled) => {
      this.setState({
        enablePPM: enabled,
      });
    });
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.NTS).then((enabled) => {
      this.setState({
        enableNTS: enabled,
      });
    });
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.NTSR).then((enabled) => {
      this.setState({
        enableNTSR: enabled,
      });
    });
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.BOAT).then((enabled) => {
      this.setState({
        enableBoat: enabled,
      });
    });
  }

  setShipmentType = (e) => {
    this.setState({ shipmentType: e.target.value });
  };

  toggleStorageModal = () => {
    this.setState((state) => ({
      showStorageInfoModal: !state.showStorageInfoModal,
    }));
  };

  toggleMoveInfoModal = () => {
    this.setState((state) => ({
      showMoveInfoModal: !state.showMoveInfoModal,
    }));
  };

  toggleBoatInfoModal = () => {
    this.setState((state) => ({
      showBoatInfoModal: !state.showBoatInfoModal,
    }));
  };

  handleSubmit = () => {
    const {
      router: { navigate },
      move,
    } = this.props;
    const { shipmentType } = this.state;

    const createShipmentPath = generatePath(customerRoutes.SHIPMENT_CREATE_PATH, { moveId: move.id });
    return navigate(`${createShipmentPath}?type=${shipmentType}`);
  };

  render() {
    const {
      router: { navigate },
      move,
      mtoShipments,
    } = this.props;
    const {
      shipmentType,
      showStorageInfoModal,
      showMoveInfoModal,
      showBoatInfoModal,
      enablePPM,
      enableNTS,
      enableNTSR,
      enableBoat,
      errorMessage,
    } = this.state;

    const shipmentInfo = determineShipmentInfo(move, mtoShipments);

    const ppmCardText = shipmentInfo.isPPMSelectable
      ? 'You pack and move your personal property or make other arrangements, The government pays you for the weight you move. This is a Personally Procured Move (PPM).'
      : 'You’ve already requested a PPM shipment. If you have more things to move yourself but that you can’t add to that shipment, contact the PPPO at your origin duty location.';

    const hhgCardText = shipmentInfo.isHHGSelectable
      ? 'All your personal property are packed and moved by professionals, paid for by the government. This is a Household Goods move (HHG).'
      : 'Talk with your movers directly if you want to add or change shipments.';

    const ntsCardText = shipmentInfo.isNTSSelectable
      ? `Movers pack and ship personal property to a storage facility, where they stay until a future move. This is an NTS (non-temporary storage) shipment.`
      : 'You’ve already requested a long-term storage shipment for this move. Talk to your movers to change or add to your request.';

    const ntsrCardText = shipmentInfo.isNTSRSelectable
      ? 'Movers pick up personal property you put into NTS during an earlier move and ship them to your new destination. This is an NTS-release (non-temporary storage release) shipment.'
      : 'You’ve already asked to have things taken out of storage for this move. Talk to your movers to change or add to your request.';

    const boatCardText = 'Provide information about your boat and we will determine how it will ship.';

    const selectableCardDefaultProps = {
      onChange: (e) => this.setShipmentType(e),
      name: 'shipmentType',
    };

    const handleBack = () => {
      const backPath = shipmentInfo.hasShipment
        ? generalRoutes.HOME_PATH
        : generatePath(customerRoutes.SHIPMENT_MOVING_INFO_PATH, { moveId: move.id });

      navigate(backPath);
    };

    return (
      <>
        <GridContainer>
          <NotificationScrollToTop dependency={errorMessage} />

          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              {errorMessage && (
                <Alert type="error" headingLevel="h4" heading="An error occurred">
                  {errorMessage}
                </Alert>
              )}

              <h6 data-testid="number-eyebrow" className="sm-heading margin-top-205 margin-bottom-0">
                Shipment {shipmentInfo.shipmentNumber}
              </h6>

              <h1 className={`${styles.selectTypeHeader} ${styles.header}`} data-testid="select-move-type-header">
                How should this shipment move?
              </h1>

              <p>
                You can move everything in one shipment, or you can split your belongings into multiple shipments that
                are moved in different ways.
              </p>
              <p>
                After you set up this shipment, you can add another shipment if you have more personal property to move.
              </p>

              <SelectableCard
                {...selectableCardDefaultProps}
                label="Movers pack and ship it, paid by the government (HHG)"
                value={SHIPMENT_OPTIONS.HHG}
                id={SHIPMENT_OPTIONS.HHG}
                cardText={hhgCardText}
                checked={shipmentType === SHIPMENT_OPTIONS.HHG}
                disabled={!shipmentInfo.isHHGSelectable}
                onHelpClick={this.toggleMoveInfoModal}
              />

              {enablePPM && (
                <SelectableCard
                  {...selectableCardDefaultProps}
                  label="Move it yourself and get paid for it (PPM)"
                  value={SHIPMENT_OPTIONS.PPM}
                  id={SHIPMENT_OPTIONS.PPM}
                  cardText={ppmCardText}
                  checked={shipmentType === SHIPMENT_OPTIONS.PPM}
                  disabled={!shipmentInfo.isPPMSelectable}
                  onHelpClick={this.toggleMoveInfoModal}
                />
              )}

              {enableNTS || enableNTSR ? (
                <>
                  <h3 className={styles.longTermStorageHeader} data-testid="long-term-storage-heading">
                    Long-term storage
                  </h3>

                  {!shipmentInfo.isNTSSelectable && !shipmentInfo.isNTSRSelectable ? (
                    <p className={styles.pSmall}>
                      Talk to your movers about long-term storage if you need to add it to this move or change a request
                      you made earlier.
                    </p>
                  ) : (
                    <>
                      <p>Your orders might not authorize long-term storage &mdash; your counselor can verify.</p>
                      {enableNTS && (
                        <SelectableCard
                          {...selectableCardDefaultProps}
                          label="It is going into storage for months or years (NTS)"
                          value={SHIPMENT_OPTIONS.NTS}
                          id={SHIPMENT_OPTIONS.NTS}
                          cardText={ntsCardText}
                          checked={shipmentType === SHIPMENT_OPTIONS.NTS && shipmentInfo.isNTSSelectable}
                          disabled={!shipmentInfo.isNTSSelectable}
                          onHelpClick={this.toggleStorageModal}
                        />
                      )}
                      {enableNTSR && (
                        <SelectableCard
                          {...selectableCardDefaultProps}
                          label="It was stored during a previous move (NTS-release)"
                          value={SHIPMENT_OPTIONS.NTSR}
                          id={SHIPMENT_OPTIONS.NTSR}
                          cardText={ntsrCardText}
                          checked={shipmentType === SHIPMENT_OPTIONS.NTSR && shipmentInfo.isNTSRSelectable}
                          disabled={!shipmentInfo.isNTSRSelectable}
                          onHelpClick={this.toggleStorageModal}
                        />
                      )}
                    </>
                  )}
                </>
              ) : null}
              {enableBoat && (
                <>
                  <h3 className={styles.longTermStorageHeader} data-testid="long-term-storage-heading">
                    Boats & Mobile Homes
                  </h3>
                  <p>
                    Moving a boat or mobile home? Please provide additional info to determine how it will be shipped.
                  </p>
                  <SelectableCard
                    {...selectableCardDefaultProps}
                    label="Move a Boat"
                    value={SHIPMENT_OPTIONS.BOAT}
                    id={SHIPMENT_OPTIONS.BOAT}
                    cardText={boatCardText}
                    checked={shipmentType === SHIPMENT_OPTIONS.BOAT && shipmentInfo.isBoatSelectable}
                    disabled={!shipmentInfo.isBoatSelectable}
                    onHelpClick={this.toggleBoatInfoModal}
                  />
                </>
              )}

              {!shipmentInfo.hasShipment && (
                <p data-testid="helper-footer" className={styles.footer}>
                  <small>
                    It is okay if you are not sure about your choice. Your move counselor will go over all your options
                    and can help make changes if necessary.
                  </small>
                </p>
              )}

              <div className={formStyles.formActions}>
                <WizardNavigation
                  disableNext={shipmentType === undefined || shipmentType === ''}
                  onBackClick={handleBack}
                  onNextClick={this.handleSubmit}
                />
              </div>
            </Grid>
          </Grid>
        </GridContainer>
        <ConnectedMoveInfoModal
          isOpen={showMoveInfoModal}
          enablePPM={enablePPM}
          closeModal={this.toggleMoveInfoModal}
        />
        <ConnectedStorageInfoModal
          isOpen={showStorageInfoModal}
          enableNTS={enableNTS}
          enableNTSR={enableNTSR}
          closeModal={this.toggleStorageModal}
        />
        <ConnectedBoatInfoModal
          isOpen={showBoatInfoModal}
          enablePPM={enableBoat}
          closeModal={this.toggleBoatInfoModal}
        />
      </>
    );
  }
}

SelectShipmentType.propTypes = {
  updateMove: func.isRequired,
  loadMTOShipments: func.isRequired,
  move: MoveTaskOrderShape.isRequired,
  mtoShipments: arrayOf(ShipmentShape).isRequired,
  router: RouterShape.isRequired,
};

const mapStateToProps = (state, ownProps) => {
  const {
    router: {
      params: { moveId },
    },
  } = ownProps;
  const move = selectMove(state, moveId);
  const mtoShipments = selectMTOShipmentsForCurrentMove(state);

  return {
    move,
    mtoShipments,
  };
};

const mapDispatchToProps = {
  updateMove: updateMoveAction,
  loadMTOShipments: loadMTOShipmentsAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(SelectShipmentType));
export { mapStateToProps as _mapStateToProps };
