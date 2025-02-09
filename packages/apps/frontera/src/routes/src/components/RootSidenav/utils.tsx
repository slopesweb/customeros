import React from 'react';

import { Triage } from '@ui/media/icons/Triage';
import { Users01 } from '@ui/media/icons/Users01.tsx';
import { HeartHand } from '@ui/media/icons/HeartHand';
import { Building07 } from '@ui/media/icons/Building07';
import { CheckHeart } from '@ui/media/icons/CheckHeart';
import { Signature } from '@ui/media/icons/Signature.tsx';
import { Briefcase01 } from '@ui/media/icons/Briefcase01';
import { BrokenHeart } from '@ui/media/icons/BrokenHeart';
import { InvoiceCheck } from '@ui/media/icons/InvoiceCheck';
import { InvoiceUpcoming } from '@ui/media/icons/InvoiceUpcoming';
import { CoinsStacked01 } from '@ui/media/icons/CoinsStacked01.tsx';
import { ClockFastForward } from '@ui/media/icons/ClockFastForward';
import { SwitchHorizontal01 } from '@ui/media/icons/SwitchHorizontal01';
export const iconMap: Record<
  string,
  (props: React.SVGAttributes<SVGElement>) => JSX.Element
> = {
  InvoiceUpcoming: (props) => <InvoiceUpcoming {...props} />,
  InvoiceCheck: (props) => <InvoiceCheck {...props} />,
  ClockFastForward: (props) => <ClockFastForward {...props} />,
  Briefcase01: (props) => <Briefcase01 {...props} />,
  Building07: (props) => <Building07 {...props} />,
  CheckHeart: (props) => <CheckHeart {...props} />,
  Seed: (props) => <HeartHand {...props} />,
  HeartHand: (props) => <HeartHand {...props} />,
  Triage: (props) => <Triage {...props} />,
  SwitchHorizontal01: (props) => <SwitchHorizontal01 {...props} />,
  BrokenHeart: (props) => <BrokenHeart {...props} />,
  CoinsStacked01: (props) => <CoinsStacked01 {...props} />,
  Signature: (props) => <Signature {...props} />,
  users_01: (props) => <Users01 {...props} />,
};
