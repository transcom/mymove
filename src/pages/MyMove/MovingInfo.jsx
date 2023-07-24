import React, { Component } from 'react';
import { connect } from 'react-redux';
import { Grid, GridContainer } from '@trussworks/react-uswds';
import { func, node, number, string } from 'prop-types';
import { generatePath } from 'react-router-dom';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './MovingInfo.module.scss';

import { customerRoutes, generalRoutes } from 'constants/routes';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { fetchLatestOrders as fetchLatestOrdersAction } from 'shared/Entities/modules/orders';
import { formatWeight } from 'utils/formatters';
import { selectCurrentOrders, selectServiceMemberFromLoggedInUser } from 'store/entities/selectors';
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
  }

  render() {
    const {
      entitlementWeight,
      router: {
        navigate,
        params: { moveId },
      },
    } = this.props;

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
              <IconSection icon="pencil-alt" headline="You don't need to get the details perfect.">
                <p>
                  After you submit this information, you&apos;ll talk to a move counselor. They will verify your choices
                  and help identify more complicated situations.
                </p>
                <p>
                  If you use movers, they will be your point of contact throughout your move and can also help you make
                  changes to your shipments.
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
              <IconSection icon="car" headline="You still have the option to move some of your belongings yourself.">
                <p>
                  Most people utilize a professional moving company to pack, pick-up and deliver the majority of their
                  personal property and move a few important or necessary items themselves. This is called a partial
                  Personally Procured Move (PPM).
                </p>
              </IconSection>
              <IconSection
                icon="hand-holding-usd"
                headline="You can get paid for any household goods you move yourself."
              >
                <p>
                  Remember to obtain and submit documents to the government to verify the weight of your PPM shipment in
                  order to receive your payment.
                </p>
              </IconSection>
            </SectionWrapper>

            <WizardNavigation
              isFirstPage
              showFinishLater
              onNextClick={() => {
                navigate(nextPath);
              }}
              onCancelClick={() => {
                navigate(generalRoutes.HOME_PATH);
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
  const entitlementWeight = orders.has_dependents
    ? serviceMember?.weight_allotment?.total_weight_self_plus_dependents
    : serviceMember?.weight_allotment?.total_weight_self;
  const serviceMemberId = serviceMember?.id;

  return {
    entitlementWeight,
    serviceMemberId,
  };
}

const mapDispatchToProps = {
  fetchLatestOrders: fetchLatestOrdersAction,
};

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(MovingInfo));
