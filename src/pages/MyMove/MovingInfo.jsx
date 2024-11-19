import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { func, node, number, string } from 'prop-types';
import { generatePath } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faSuitcase } from '@fortawesome/free-solid-svg-icons';

import { isBooleanFlagEnabled } from '../../utils/featureFlags';
import { FEATURE_FLAG_KEYS } from '../../shared/constants';

import styles from './MovingInfo.module.scss';

import { customerRoutes, generalRoutes } from 'constants/routes';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { fetchLatestOrders as fetchLatestOrdersAction } from 'shared/Entities/modules/orders';
import { formatUBAllowanceWeight, formatWeight } from 'utils/formatters';
import { selectCurrentOrders, selectServiceMemberFromLoggedInUser, selectUbAllowance } from 'store/entities/selectors';
import withRouter from 'utils/routing';
import { RouterShape } from 'types';

const IconSection = ({ icon, headline, children }) => (
  <Grid row className={styles.IconSection}>
    <Grid col="auto" className={styles.SectionIcon}>
      <FontAwesomeIcon size="lg" icon={icon} />
    </Grid>
    <Grid col="fill" className={styles.SectionContent}>
      <h2 className={styles.SectionHeadline}>{headline} </h2>
      {children}
    </Grid>
  </Grid>
);

IconSection.propTypes = {
  icon: string.isRequired,
  headline: string.isRequired,
  children: node.isRequired,
};

export class MovingInfo extends Component {
  componentDidMount() {
    const { serviceMemberId, fetchLatestOrders } = this.props;
    fetchLatestOrders(serviceMemberId);
    isBooleanFlagEnabled('multi_move').then((enabled) => {
      this.setState({
        multiMoveFeatureFlag: enabled,
      });
    });
    isBooleanFlagEnabled(FEATURE_FLAG_KEYS.PPM).then((enabled) => {
      this.setState({
        ppmFeatureFlag: enabled,
      });
    });
  }

  render() {
    const {
      ubAllowance,
      entitlementWeight,
      router: {
        navigate,
        params: { moveId },
      },
    } = this.props;

    let multiMove = false;
    let enablePPM = true;
    if (this.state) {
      const { multiMoveFeatureFlag, ppmFeatureFlag } = this.state;
      multiMove = multiMoveFeatureFlag;
      enablePPM = ppmFeatureFlag;
    }

    const nextPath = generatePath(customerRoutes.SHIPMENT_SELECT_TYPE_PATH, {
      moveId,
    });

    return (
      <GridContainer>
        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <h1 className={styles.ShipmentsHeader}>Things to know about selecting shipments</h1>
            <SectionWrapper className={styles.Wrapper}>
              <IconSection
                icon="weight-hanging"
                headline={`You can move ${formatWeight(entitlementWeight)} in this move.`}
              >
                <p>
                  You will have to pay for any excess weight above this allowance, so work hard to make sure you stay
                  within your weight limit.
                </p>
              </IconSection>
              {ubAllowance ? (
                <IconSection
                  icon={faSuitcase}
                  headline={`You can move up to ${formatUBAllowanceWeight(
                    ubAllowance,
                  )} of unaccompanied baggage in this move.`}
                >
                  <p>
                    If you request an unaccompanied baggage (UB) shipment, keep in mind that you will need to stay under
                    that weight allowance for your UB shipment, and that the weight of your UB shipment is also part of
                    your overall authorized weight allowance.
                  </p>
                </IconSection>
              ) : null}
              <IconSection icon="pencil-alt" headline="You don't need to get the details perfect.">
                <p>
                  After you submit this information, you will talk to a government move counselor. They will verify your
                  choices and help identify more complicated situations.
                </p>
                <p>
                  When counseling is complete and a Move Task Order (MTO) is issued, you will be appointed a Customer
                  Care Representative. They will be your point of contact for the rest of your move and can help make
                  any changes to your shipment.
                </p>
              </IconSection>
              <IconSection icon="truck-moving" headline="Your Move Manager will:">
                <div className={styles.IconSectionList}>
                  <ul>
                    <li>Help estimate how much your belongings weigh</li>
                    <li>Set pack and pickup dates based on your preferred pickup date</li>
                    <li>Be your main point of contact during your move</li>
                  </ul>
                </div>
              </IconSection>
              {enablePPM ? (
                <>
                  <IconSection
                    icon="car"
                    headline="You still have the option to move some of your belongings yourself."
                  >
                    <p>
                      Most people utilize a professional moving company to pack, pick-up and deliver the majority of
                      their personal property and move a few important or necessary items themselves. This is called a
                      partial Personally Procured Move (PPM).
                    </p>
                  </IconSection>
                  <IconSection
                    icon="hand-holding-usd"
                    headline="You can get paid for any household goods you move yourself."
                  >
                    <p>
                      Remember to obtain and submit documents to the government to verify the weight of your PPM
                      shipment in order to receive your payment.
                    </p>
                  </IconSection>
                </>
              ) : null}
            </SectionWrapper>

            <WizardNavigation
              isFirstPage
              showFinishLater
              onNextClick={() => {
                navigate(nextPath);
              }}
              onCancelClick={() => {
                if (multiMove) {
                  navigate(generatePath(customerRoutes.MOVE_HOME_PATH, { moveId }));
                } else {
                  navigate(generalRoutes.HOME_PATH);
                }
              }}
            />
          </Grid>
        </Grid>
      </GridContainer>
    );
  }
}

MovingInfo.propTypes = {
  entitlementWeight: number,
  fetchLatestOrders: func.isRequired,
  serviceMemberId: string.isRequired,
  router: RouterShape,
};

MovingInfo.defaultProps = {
  entitlementWeight: 0,
  router: {},
};

function mapStateToProps(state) {
  const orders = selectCurrentOrders(state);
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const entitlementWeight = orders.authorizedWeight;
  const serviceMemberId = serviceMember?.id;
  const ubAllowance = selectUbAllowance(state);

  return {
    ubAllowance,
    entitlementWeight,
    serviceMemberId,
  };
}

const mapDispatchToProps = {
  fetchLatestOrders: fetchLatestOrdersAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MovingInfo));
