(MessageG ((label Suggest) (payload ((PValue () PTString)))) B A
 (MessageG ((label Query) (payload ((PValue () PTString)))) A S
  (ChoiceG S
   ((MessageG ((label Available) (payload ((PValue () (PTAbstract number)))))
     S A
     (MessageG ((label Quote) (payload ((PValue () (PTAbstract number))))) A
      B
      (ChoiceG B
       ((MessageG ((label OK) (payload ((PValue () (PTAbstract number))))) B
         A
         (MessageG
          ((label Confirm) (payload ((PValue () (PTAbstract Cred))))) A S
          EndG))
        (MessageG ((label No) (payload ())) B A
         (MessageG ((label Reject) (payload ())) A S EndG))))))
    (MessageG ((label Full) (payload ())) S A
     (MessageG ((label Full) (payload ())) A B EndG))))))
