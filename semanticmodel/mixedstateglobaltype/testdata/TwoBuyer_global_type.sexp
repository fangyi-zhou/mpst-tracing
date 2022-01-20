(MessageG ((label query) (payload ())) A S
 (MessageG ((label quote) (payload ())) S A
  (MessageG ((label quote) (payload ())) S B
   (MessageG ((label share) (payload ())) B A
    (ChoiceG A
     ((MessageG ((label buy) (payload ())) A S EndG)
      (MessageG ((label cancel) (payload ())) A S EndG)))))))
