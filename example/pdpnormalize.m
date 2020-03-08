x=[1 .5 2 1.5 3]
tt=[0 0.001 0.00210 0.00320 0.015]
ts=min(diff(tt));

stem(tt,x,'r*');
newtt=0:ts:30*ts
grid on;
hold all;
newpdp=zeros(1,length(newtt));
for k=1:length(x)
plot(newtt,x(k)*sinc((newtt-tt(k))/ts))
newpdp=newpdp+x(k)*sinc((newtt-tt(k))/ts);
end
stem(newtt,newpdp,'b--')
xlabel('\tau (s)')
ylabel('Power')
title('PDP and resampling')
legend('original PDP','normalized ts=1e-3')
