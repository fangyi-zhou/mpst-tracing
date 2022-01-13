(ChoiceG M
 ((MessageG ((label STOP_from_A) (payload ())) A M
   (MessageG ((label STOPPED_A) (payload ())) M A
    (MessageG ((label STOPPED_A) (payload ())) M B EndG)))
  (MessageG ((label STOP_from_B) (payload ())) B M
   (MessageG ((label STOPPED_B) (payload ())) M A
    (MessageG ((label STOPPED_B) (payload ())) M B EndG)))))
